package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/webhook"
)

type serveFlags struct {
	port      int
	secretEnv string
	secret    string
}

func newServeCmd() *cobra.Command {
	var f serveFlags
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start an HTTP server that receives and verifies Fireflies webhook events (V1 and V2)",
		Long: `Starts an HTTP server that listens for Fireflies webhook POST requests.

Verified events are emitted as NDJSON on stdout.
Rejected events (bad or missing signature) return 401 and log a warning to stderr.

Routes:
  POST /webhooks/v1  — Fireflies webhook V1
  POST /webhooks/v2  — Fireflies webhook V2
  GET  /health       — Health check (200 ok)

The webhook secret must be set via --secret-env (env var name) or --secret.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(f)
		},
	}
	cmd.Flags().IntVar(&f.port, "port", 8080, "Port to listen on")
	cmd.Flags().StringVar(&f.secretEnv, "secret-env", "", "Name of env var holding the webhook secret")
	cmd.Flags().StringVar(&f.secret, "secret", "", "Webhook secret value (use --secret-env in production)")
	return cmd
}

func runServe(f serveFlags) error {
	// Resolve secret.
	secret := f.secret
	if f.secretEnv != "" {
		secret = os.Getenv(f.secretEnv)
	}
	if secret == "" {
		return fmt.Errorf("webhook secret is empty: set --secret-env <VAR> or --secret <value>")
	}

	mux := buildMux([]byte(secret))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", f.port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM.
	idleConnsClosed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Printf("[webhooks] shutting down (signal)")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("[webhooks] shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("[webhooks] listening on :%d", f.port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	<-idleConnsClosed
	return nil
}

// buildMux constructs the HTTP mux for testing without starting a server.
func buildMux(secret []byte) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/webhooks/v1", makeWebhookHandler(secret, "v1"))
	mux.HandleFunc("/webhooks/v2", makeWebhookHandler(secret, "v2"))
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "ok")
}

func makeWebhookHandler(secret []byte, version string) http.HandlerFunc {
	// Structured access logger — logs to stderr, redacts body/sig.
	logger := log.New(os.Stderr, "[webhooks] ", log.LstdFlags)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			logger.Printf("method=%s path=%s status=400 dur=%s err=%q", r.Method, r.URL.Path, time.Since(start), err)
			return
		}
		defer r.Body.Close() //nolint:errcheck // request body already fully read above

		sigHeader := r.Header.Get("x-hub-signature")

		if err := webhook.Verify(secret, body, sigHeader); err != nil {
			// Log warning WITHOUT the signature or body to avoid leaking secrets.
			logger.Printf("method=%s path=%s status=401 dur=%s msg=%q", r.Method, r.URL.Path, time.Since(start), "signature verification failed")
			http.Error(w, "unauthorized: signature verification failed", http.StatusUnauthorized)
			return
		}

		// Decode into generic map.
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			// Treat malformed JSON as a bad request but still accept it for
			// logging (signature was valid).
			payload = map[string]any{"raw": string(body)}
		}
		payload["version"] = version

		// Emit one NDJSON line on stdout.
		line, err := json.Marshal(payload)
		if err == nil {
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", line)
		}

		w.WriteHeader(http.StatusOK)
		logger.Printf("method=%s path=%s status=200 dur=%s", r.Method, r.URL.Path, time.Since(start))
	}
}
