package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Connection struct {
	*websocket.Conn
	PublicKey *rsa.PublicKey
}

// Consrtructs a new Connection object
func NewConnection(conn *websocket.Conn, publicKey *rsa.PublicKey) *Connection {
	return &Connection{Conn: conn, PublicKey: publicKey}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var connections = make(map[string]*websocket.Conn) // Map of player IDs to connections

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	playerID := r.URL.Query().Get("playerID")
	connections[playerID] = conn

	for {
		// Read messages from the player
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			delete(connections, playerID)
			break
		}

		fmt.Printf("MESSAGE (%v): %v\n", playerID, string(message))

		// Broadcast to other players
		for id, otherConn := range connections {
			if id != playerID {
				otherConn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}
