package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type connectionManager struct {
	mu                 sync.Mutex
	GameToPlayers      map[string][]string
	PlayerToGame       map[string]string
	PlayerToConnection map[string]*websocket.Conn
}

func NewConnectionManager() *connectionManager {
	gameToPlayers := make(map[string][]string)
	playerToGame := make(map[string]string)
	playerToConnection := make(map[string]*websocket.Conn)

	return &connectionManager{
		GameToPlayers:      gameToPlayers,
		PlayerToGame:       playerToGame,
		PlayerToConnection: playerToConnection,
	}
}

// Returns the GameIDs for all games in the connection manager
func (cm *connectionManager) Games() []string {
	games := make([]string, 0)
	for key, _ := range cm.GameToPlayers {
		games = append(games, key)
	}

	return games
}

// Returns all PlayerIDs for players with an active connection in
// the connection manager
func (cm *connectionManager) Players() []string {
	players := make([]string, 0)
	for key, _ := range cm.PlayerToConnection {
		players = append(players, key)
	}

	return players
}

func (cm *connectionManager) Add(gameId, playerId string, connection *websocket.Conn) error {
	// Locking connectionManager to ensure data is not corrupted by multiple
	// writes
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Checking if the game already exists
	if !slices.Contains(cm.Games(), gameId) {
		return &GameNotInMapError{}
	}

	// Check if player is already connected to game
	if _, ok := cm.PlayerToGame[playerId]; ok {
		return &PlayerAlreadyConnectedError{}
	}

	cm.GameToPlayers[gameId] = append(cm.GameToPlayers[gameId], playerId)
	cm.PlayerToGame[playerId] = gameId
	cm.PlayerToConnection[playerId] = connection

	return nil
}

func (cm *connectionManager) Remove(playerId string) error {
	// Locking connectionManager to ensure data is not corrupted by multiple
	// writes
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Checking if player is connected
	// if not we can return without modifications
	if _, ok := cm.PlayerToConnection[playerId]; !ok {
		return nil
	}

	// Removes player from GameToPlayers
	// by finding it in list and removing it
	game := cm.PlayerToGame[playerId]
	playersInGame := cm.GameToPlayers[game]
	idx := -1
	for ix, playerInGame := range playersInGame {
		if playerId == playerInGame {
			idx = ix
		}
	}
	cm.GameToPlayers[game] = slices.Delete(cm.GameToPlayers[game], idx, idx+1)

	// Removes player in PlayerToGame
	delete(cm.PlayerToGame, playerId)

	// Removes player in PlayerToConnection
	delete(cm.PlayerToConnection, playerId)

	return nil
}

func (cm *connectionManager) Broadcast(message []byte, playerIds []string) {
	for _, playerId := range playerIds {
		connection := cm.PlayerToConnection[playerId]
		connection.WriteMessage(websocket.TextMessage, message)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var connections = NewConnectionManager()

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
	gameId := r.URL.Query().Get("game_id")
	if gameId == "" {
		// Make new lobby logic
	} else {
		connections.Add(gameId, playerId, conn)
	}

	for {
		// Read messages from the player
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			connections.Remove(playerId)
			break
		}

		fmt.Printf("MESSAGE (%v): %v\n", playerId, string(message))

		// Broadcasting message to every player except current one
		connections.Broadcast(message, connections.Players())
	}
}
