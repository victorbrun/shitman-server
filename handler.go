package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
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

// var connections = make(map[string]*websocket.Conn)
var games = make(map[string]*Game)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Extracts player ID
	playerId := r.URL.Query().Get("player_id")
	if playerId == "" {
		log.Printf("No player ID specified. Closing connection")
		conn.WriteJSON(http.StatusBadRequest)
		conn.Close()
		return
	}

	// Creates game if not exists
	gameId := r.URL.Query().Get("game_id")
	if gameId == "" {
		gameId = uuid.New().String()
		games[gameId] = NewGame(gameId)
		log.Printf("Created game: %v by player: %v", gameId, playerId)
	} else if _, ok := games[gameId]; !ok {
		log.Printf("Game: %v does not exist. Closing connection", gameId)
		conn.WriteJSON(http.StatusNotFound)
		conn.Close()
		return
	}

	// Add player to game
	games[gameId].AddPlayer(playerId)

	conn.WriteJSON(http.StatusOK)
	conn.Close()
}
