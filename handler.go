package main

import (
	"crypto/rsa"
	"encoding/base64"
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

var connections = make(map[string]*Connection) // Map of player IDs to connections

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Extracts auth info from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {

	}

	// Extracts public key from header
	publicKeyBase64 := r.Header.Get("X-Public-Key")
	if publicKeyBase64 == "" {
		http.Error(w, "Missing public key", http.StatusBadRequest)
		return
	}

	// Decodes public key
	publicKeyBites, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		http.Error(w, "Invalid public key encoding", http.StatusBadRequest)
		return
	}

	// Constructs PublicKey object from decoded key
	publicKey, err := parsePublicKeyPEM(publicKeyBites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	playerID := r.URL.Query().Get("playerID")
	connections[playerID] = NewConnection(conn, publicKey)

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
