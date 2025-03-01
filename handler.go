package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type PlayerCommand struct {
	// Metadata
	ID string `json:"player_id"`

	// List of cards player wants to play.
	// The order of the list defines the order
	// in which the cards ought to be played
	PlayCards []Card `json:"play_cards"`

	// If card from hidden hand should be played
	// as last action
	PlayCardFromHiddenHand bool `json:"play_card_from_hidden_hand"`

	// If top most card from deck ought to be played.
	// This can only be used when no other card can be played
	// by the player.
	PlayRandomCardFromDeck bool `json:"play_random_card_from_deck"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var connections = make(map[string]*websocket.Conn)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Extracts player ID and game id and creats new lobby if
	// no id is given
	playerId := r.URL.Query().Get("player_id")
	//gameId := r.URL.Query().Get("game_id")
	connections[playerId] = conn

	for {
		// Read messages from the player
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
		}

		// Parsing player command
		var command PlayerCommand
		err = json.Unmarshal(message, &command)
		if err != nil {
			log.Println("Error parsing message:", err)
		}

		log.Printf("Parsed command (%v): %+v\n", playerId, command)

		conn.WriteJSON(http.StatusOK)
	}
}
