package realtime

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// TestSubscribe_DispatchEvents spins up a minimal httptest WS server that
// sends Socket.IO-framed packets and verifies that Subscribe dispatches the
// correct Event values to the handler.
func TestSubscribe_DispatchEvents(t *testing.T) {
	type serverMsg struct {
		name    string
		payload any
	}

	want := []serverMsg{
		{name: "transcript", payload: map[string]any{"text": "hello"}},
		{name: "speaker", payload: map[string]any{"speaker": "Alice"}},
	}

	var received []Event
	var mu sync.Mutex
	done := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer conn.Close() //nolint:errcheck

		// Send Engine.IO open packet.
		if err := conn.WriteMessage(websocket.TextMessage, []byte(`0{"sid":"abc","upgrades":[],"pingInterval":25000,"pingTimeout":5000}`)); err != nil {
			return
		}

		// Read the CONNECT packet from client.
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if string(msg) != "40" {
			t.Errorf("expected Socket.IO CONNECT '40', got %q", msg)
		}

		// Read the subscribe event from client.
		_, msg, err = conn.ReadMessage()
		if err != nil {
			return
		}
		if !strings.HasPrefix(string(msg), `42["subscribe"`) {
			t.Errorf("expected subscribe event, got %q", msg)
		}

		// Send a ping and verify pong.
		if err := conn.WriteMessage(websocket.TextMessage, []byte("2")); err != nil {
			return
		}
		_, pong, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if string(pong) != "3" {
			t.Errorf("expected pong '3', got %q", pong)
		}

		// Send Socket.IO EVENT packets.
		for _, m := range want {
			payloadBytes, _ := json.Marshal(m.payload)
			pkt := `42["` + m.name + `",` + string(payloadBytes) + `]`
			if err := conn.WriteMessage(websocket.TextMessage, []byte(pkt)); err != nil {
				return
			}
		}

		// Wait for client to close (context cancelled).
		conn.ReadMessage() //nolint:errcheck
	}))
	defer server.Close()

	// Convert http:// → ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	c := New("testkey", wsURL)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Subscribe(ctx, "meeting-123", func(e Event) { //nolint:errcheck
			mu.Lock()
			received = append(received, e)
			if len(received) == len(want) {
				close(done)
			}
			mu.Unlock()
		})
	}()

	select {
	case <-done:
		cancel()
	case <-ctx.Done():
		t.Error("timed out waiting for events")
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(received) != len(want) {
		t.Fatalf("received %d events, want %d", len(received), len(want))
	}
	for i, w := range want {
		if received[i].Name != w.name {
			t.Errorf("[%d] name: got %q, want %q", i, received[i].Name, w.name)
		}
	}
}

// TestNew_DefaultEndpoint verifies that an empty endpoint uses DefaultEndpoint.
func TestNew_DefaultEndpoint(t *testing.T) {
	c := New("key", "")
	if c.endpoint != DefaultEndpoint {
		t.Errorf("endpoint: got %q, want %q", c.endpoint, DefaultEndpoint)
	}
}
