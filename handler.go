package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

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
		conn.WriteJSON(map[string]string{"error": "No player ID specified"})
		conn.Close()
		return
	}

	// Creates game if not exists
	gameId := r.URL.Query().Get("game_id")
	if gameId == "" {
		gameId = uuid.New().String()
		games[gameId] = NewGame(gameId)

		log.Printf("Created game: %v by player: %v", gameId, playerId)
		conn.WriteJSON(map[string]string{"message": "Created game successfully", "game_id": gameId})

	} else if _, ok := games[gameId]; !ok {
		log.Printf("Game: %v does not exist. Closing connection", gameId)
		conn.WriteJSON(map[string]string{"error": "Game does not exist", "game_id": gameId})
		conn.Close()
		return
	}

	// Add player to game
	err = games[gameId].AddPlayer(playerId, conn)
	if err != nil {
		log.Printf("Error adding player to game: %v", err)
		conn.WriteJSON(map[string]string{"error": "Could not add player to game", "player_id": playerId, "game_id": gameId})
		conn.Close()
	}
	log.Printf("Player: %v succesfully joined game: %v", playerId, gameId)
	conn.WriteJSON(map[string]string{"message": "Joined game successfully", "game_id": gameId})

	// Extracting game for easy access
	game := games[gameId]

	// Extracing player for easy access
	player := game.findPlayerById(playerId)

	// Start listening for commands from client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from player %s: %v", playerId, err)
			conn.WriteJSON(map[string]string{"error": "Error reading message from player", "player_id": player.ID})
		}

		// Handle the message
		err = game.handlePlayerMessage(player, message)
		if err == nil {
			continue
		}

		// Error handling
		switch e := err.(type) {
		case *GameNotInLobbyError:
			player.conn.WriteJSON(map[string]string{"error": "Game not in lobby", "player_id": gameId})
		case *NotGameOwnerError:
			player.conn.WriteJSON(map[string]string{"error": "Player not owner of game", "player_id": player.ID, "game_id": gameId})
		case *GameNotStartedError:
			player.conn.WriteJSON(map[string]string{"error": "Game is not started", "game_id": gameId})
		case *InvalidArgumentError:
			player.conn.WriteJSON(map[string]string{"error": "Invalid argument"})
		case *NotPlayersTurnError:
			player.conn.WriteJSON(map[string]string{"error": "Not players turn", "player_id": player.ID})
		case *CardCannotBePlayedError:
			player.conn.WriteJSON(map[string]string{
				"error":                 "Card cannot be played",
				"card_to_play":          e.CardToPlay.String(),
				"card_on_playing_field": e.CardOnPlayingField.String(),
			})
		case *CardAlreadyPlayedError:
			player.conn.WriteJSON(map[string]string{"error": "Card has already been played"})
		default:
			player.conn.WriteJSON(map[string]string{"error": "Internal server error"})
		}
	}

}
