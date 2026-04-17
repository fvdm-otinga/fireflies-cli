// Package realtime provides a Socket.IO v4 WebSocket client for the
// Fireflies realtime transcript event stream at wss://api.fireflies.ai/ws/realtime.
package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const DefaultEndpoint = "wss://api.fireflies.ai/ws/realtime"

// Event is a decoded Socket.IO event delivered to the handler.
type Event struct {
	Name    string
	Payload json.RawMessage
}

// Client connects to the Fireflies Socket.IO realtime endpoint.
type Client struct {
	apiKey   string
	endpoint string
}

// New creates a new realtime Client.
func New(apiKey, endpoint string) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	return &Client{apiKey: apiKey, endpoint: endpoint}
}

// Subscribe connects to the realtime endpoint, subscribes to events for
// meetingID, and invokes handler for each received event. It reconnects
// with exponential backoff on network errors. It blocks until ctx is
// cancelled.
func (c *Client) Subscribe(ctx context.Context, meetingID string, handler func(Event)) error {
	backoff := 500 * time.Millisecond
	const maxBackoff = 30 * time.Second

	for {
		err := c.connect(ctx, meetingID, handler)
		if ctx.Err() != nil {
			// Context cancelled — clean exit.
			return nil
		}
		if err != nil {
			// Reconnect with exponential backoff.
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}
		return nil
	}
}

// connect establishes one WS session and reads events until the connection
// drops or ctx is cancelled.
func (c *Client) connect(ctx context.Context, meetingID string, handler func(Event)) error {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+c.apiKey)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.endpoint, headers)
	if err != nil {
		return fmt.Errorf("websocket dial: %w", err)
	}
	defer conn.Close()

	// Close the WS connection when ctx is done.
	go func() {
		<-ctx.Done()
		conn.WriteMessage( //nolint:errcheck
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		conn.Close()
	}()

	// Socket.IO v4 framing:
	//   "0{...}"  — Engine.IO open (server sends first)
	//   "40"       — Socket.IO CONNECT on default namespace
	//   "42[...]"  — Socket.IO EVENT
	//   "2"        — Engine.IO ping
	//   "3"        — Engine.IO pong (client reply)

	// Wait for Engine.IO open packet.
	if err := waitForOpen(conn); err != nil {
		return err
	}

	// Send Socket.IO CONNECT + subscribe payload.
	if err := sendConnect(conn, meetingID); err != nil {
		return fmt.Errorf("send connect: %w", err)
	}

	// Read loop.
	for {
		if ctx.Err() != nil {
			return nil
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("read: %w", err)
		}
		handlePacket(conn, msg, handler)
	}
}

// waitForOpen reads packets until it sees an Engine.IO "0" open packet.
func waitForOpen(conn *websocket.Conn) error {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second)) //nolint:errcheck
	defer conn.SetReadDeadline(time.Time{})                //nolint:errcheck
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("waiting for open: %w", err)
		}
		s := string(msg)
		if strings.HasPrefix(s, "0") {
			return nil
		}
		// Ignore anything else during handshake.
	}
}

// sendConnect sends the Socket.IO CONNECT packet and a subscribe event.
func sendConnect(conn *websocket.Conn, meetingID string) error {
	// Socket.IO CONNECT on default namespace.
	if err := conn.WriteMessage(websocket.TextMessage, []byte("40")); err != nil {
		return err
	}
	// Subscribe event: 42["subscribe", {"meetingId": "<id>"}]
	payload, err := json.Marshal(map[string]string{"meetingId": meetingID})
	if err != nil {
		return err
	}
	pkt := fmt.Sprintf(`42["subscribe",%s]`, string(payload))
	return conn.WriteMessage(websocket.TextMessage, []byte(pkt))
}

// handlePacket dispatches a raw Socket.IO packet.
func handlePacket(conn *websocket.Conn, msg []byte, handler func(Event)) {
	s := string(msg)
	switch {
	case s == "2":
		// Engine.IO ping — reply with pong.
		conn.WriteMessage(websocket.TextMessage, []byte("3")) //nolint:errcheck
	case strings.HasPrefix(s, "42"):
		// Socket.IO EVENT packet: 42["name", payload]
		inner := strings.TrimPrefix(s, "42")
		var parts []json.RawMessage
		if err := json.Unmarshal([]byte(inner), &parts); err != nil || len(parts) < 1 {
			return
		}
		var name string
		if err := json.Unmarshal(parts[0], &name); err != nil {
			return
		}
		var payload json.RawMessage
		if len(parts) > 1 {
			payload = parts[1]
		}
		handler(Event{Name: name, Payload: payload})
	}
	// Other packet types (open, close, ack, etc.) are silently ignored.
}
