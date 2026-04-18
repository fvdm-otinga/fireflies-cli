package webhooks

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/webhook"
)

// webhookSecretEnv is the fixed env var the command reads the secret from.
const webhookSecretEnv = "FIREFLIES_WEBHOOK_SECRET"

// maxBodyBytes caps the webhook request body at 1 MiB. Requests exceeding this
// are rejected with HTTP 413 before any HMAC work is done, mitigating
// pre-authentication DoS via oversized payloads.
const maxBodyBytes int64 = 1 << 20

type serveFlags struct {
	port        int
	secretEnv   string // optional override: read secret from this env var name
	secretStdin bool   // read one line from stdin and use as secret
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

The webhook secret is read from the FIREFLIES_WEBHOOK_SECRET env var by default,
or (with --secret-stdin) from a single line on stdin. Use --secret-env <NAME>
to read from a differently-named env var. Passing the secret as a command-line
flag is not supported (it would leak via the process listing).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(f)
		},
	}
	cmd.Flags().IntVar(&f.port, "port", 8080, "Port to listen on")
	cmd.Flags().StringVar(&f.secretEnv, "secret-env", "", "Name of env var holding the webhook secret (defaults to FIREFLIES_WEBHOOK_SECRET)")
	cmd.Flags().BoolVar(&f.secretStdin, "secret-stdin", false, "Read the webhook secret as a single line from stdin")
	return cmd
}

func runServe(f serveFlags) error {
	secret, err := resolveSecret(f, os.Stdin)
	if err != nil {
		return err
	}

	mux := buildMux([]byte(secret))
	handler := recoverMiddleware(mux)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", f.port),
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
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

// resolveSecret picks the webhook secret from (in order): --secret-stdin,
// --secret-env <NAME>, then the default FIREFLIES_WEBHOOK_SECRET env var.
// Returns a usage error (exit code 2) if no secret is available.
func resolveSecret(f serveFlags, stdin io.Reader) (string, error) {
	if f.secretStdin {
		br := bufio.NewReader(stdin)
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", ferr.Usage(fmt.Sprintf("failed to read secret from stdin: %v", err))
		}
		secret := strings.TrimSpace(line)
		if secret == "" {
			return "", ferr.Usage("webhook secret from stdin is empty")
		}
		return secret, nil
	}

	envName := webhookSecretEnv
	if f.secretEnv != "" {
		envName = f.secretEnv
	}
	if secret := os.Getenv(envName); secret != "" {
		return secret, nil
	}

	return "", ferr.Usage(fmt.Sprintf(
		"webhook secret is empty: set %s env var, or pass --secret-stdin, or use --secret-env <VAR>",
		webhookSecretEnv,
	))
}

// recoverMiddleware protects the mux from handler panics: a recovered panic is
// logged and the client receives a 500 instead of a dropped connection.
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				log.Printf("panic: %v", rec)
			}
		}()
		next.ServeHTTP(w, r)
	})
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

		// Cap request body size before any allocation/HMAC work. Exceeding the
		// limit surfaces as *http.MaxBytesError from io.ReadAll and we respond
		// 413 with a JSON error body.
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			var mbe *http.MaxBytesError
			if errors.As(err, &mbe) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				_, _ = fmt.Fprint(w, `{"error":"payload too large"}`)
				logger.Printf("method=%s path=%s status=413 dur=%s msg=%q", r.Method, r.URL.Path, time.Since(start), "payload too large")
				return
			}
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
