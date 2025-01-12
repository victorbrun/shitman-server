package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketServer(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(handleConnection))
	defer server.Close()

	// Convert test server URL to WebSocket URL
	wsURL := "ws" + server.URL[len("http:"):]

	// Number of test clients
	numClients := 3
	var wg sync.WaitGroup
	wg.Add(numClients)

	// Message to be broadcast
	testMessage := "Hello, WebSocket!"

	// Channel to collect received messages
	receivedMessages := make(chan string, numClients)

	// Spin up test clients
	for i := 1; i <= numClients; i++ {
		go func(clientID int) {
			defer wg.Done()

			// Create a WebSocket client
			dialer := websocket.Dialer{}
			conn, _, err := dialer.Dial(wsURL+"?player_id=player"+string(rune(clientID))+"&game_id=testgame", nil)
			if err != nil {
				t.Fatalf("Client %d: failed to connect: %v", clientID, err)
			}
			defer conn.Close()

			// Send a message if this is the first client
			if clientID == 1 {
				err := conn.WriteMessage(websocket.TextMessage, []byte(testMessage))
				if err != nil {
					t.Fatalf("Client %d: failed to send message: %v", clientID, err)
				}
			}

			// Receive a message
			_, message, err := conn.ReadMessage()
			if err != nil {
				t.Fatalf("Client %d: failed to read message: %v", clientID, err)
			}

			receivedMessages <- string(message)
		}(i)
	}

	// Wait for all clients to finish
	wg.Wait()
	close(receivedMessages)

	// Verify all clients received the message
	for msg := range receivedMessages {
		if msg != testMessage {
			t.Errorf("Received unexpected message: %s", msg)
		}
	}
}
