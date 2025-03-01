package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type GameStatus int

const (
	InLobby GameStatus = iota
	InGame
	Finished
	Canceled
)

type GameCommand struct {
	PlayerId string `json:"player_id"`
	GameId   string `json:"game_id"`
	Action   string `json:"game_action"`
}

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

type Player struct {
	// Metadata
	ID   string
	conn *websocket.Conn

	// Cards
	PrivateHand *Hand
	PublicHand  *Hand
	HiddenHand  *Hand
}

// Version of player that can be exposed to other players
type PublicCensoredPlayer struct {
	// Metadata
	ID string

	// Cards that can be shown for the public
	// and the number of cards in the other hands
	PrivateHandCount int
	PublicHand       *Hand
	HiddenHandCount  int
}

// Version of player that can be exposed to the player itself
type PrivateCensoredPlayer struct {
	// Metadata
	ID string

	// Cards that can be shown for the player
	// and the number of cards in the other hands
	PrivateHand     *Hand
	PublicHand      *Hand
	HiddenHandCount int
}

type PlayedCard struct {
	Card          Card      `json:"card"`
	PlayedAt      time.Time // Unix timestamp for when card was played
	PlayedBy      *Player
	PlayedInRound int
}

func NewPlayer(id string, conn *websocket.Conn, deck *Deck) *Player {
	return &Player{
		ID:   id,
		conn: conn,

		PublicHand:  NewHand(3, deck),
		PrivateHand: NewHand(3, deck),
		HiddenHand:  NewHand(3, deck),
	}
}

func NewPlayedCard(card Card, playedBy *Player, playedInRound int) PlayedCard {
	return PlayedCard{
		Card:          card,
		PlayedAt:      time.Now().UTC(),
		PlayedBy:      playedBy,
		PlayedInRound: playedInRound,
	}
}

// Player plays card from hand onto pf
func (p *Player) Play(card Card, pf *PlayingField, round int) error {
	// Checks that player can play card
	idx, ok := p.PrivateHand.Contains(card)
	if !ok {
		return &CardNotInCollectionError{}
	}

	// Constructs PlayedCard and puts it onto pf
	pc := NewPlayedCard(card, p, round)
	pf.ActivePlayedCards = append(pf.ActivePlayedCards, pc)

	// Removes card from hand
	p.PrivateHand.Remove(idx)

	// Returns without error
	return nil
}

// Censors the hand such that the information can be
// communicated to every participant in the game
func (p *Player) PublicCensor() *PublicCensoredPlayer {
	return &PublicCensoredPlayer{
		ID:               p.ID,
		PrivateHandCount: p.PrivateHand.Size(),
		PublicHand:       p.PublicHand,
		HiddenHandCount:  p.HiddenHand.Size(),
	}
}

// Censors the hands such that the information can be communicated
// to the player owning this hand
func (p *Player) PrivateCensor() *PrivateCensoredPlayer {
	return &PrivateCensoredPlayer{
		ID:              p.ID,
		PrivateHand:     p.PrivateHand,
		PublicHand:      p.PublicHand,
		HiddenHandCount: p.HiddenHand.Size(),
	}
}

type PlayingField struct {
	Deck                *Deck
	ActivePlayedCards   []PlayedCard
	InactivePlayedCards []PlayedCard
}

func NewPlayingField(deck *Deck) *PlayingField {
	return &PlayingField{
		Deck:                deck,
		ActivePlayedCards:   make([]PlayedCard, 0),
		InactivePlayedCards: make([]PlayedCard, 0),
	}
}

// Returns the top card in tne acive played cards, i.e. the card
// against which one plays
func (pf *PlayingField) Top() Card {
	return pf.ActivePlayedCards[len(pf.ActivePlayedCards)-1].Card
}

func (pf *PlayingField) TestToPlayCard(card Card) error {
	// If no active played cards all cards can be played
	if len(pf.ActivePlayedCards) == 0 {
		return nil
	}

	// If card trying to be played is 2 or 10 it can
	// always be played
	if card.Rank.Value(true) == 2 || card.Rank.Value(true) == 10 {
		return nil
	}

	// Checking skitgubbe rules against top card
	if card.Rank.Value(true) >= pf.Top().Rank.Value(true) {
		return nil
	}

	// If we have not yet returned, the card cannot be played
	return &CardCannotBePlayedError{CardToPlay: card, CardOnPlayingField: pf.Top()}
}

type Game struct {
	// Metadata
	ID string

	// Player data
	Owner   *Player
	Players []*Player

	// Game data
	PlayingField     *PlayingField
	Status           GameStatus
	Round            int
	PlayOrder        []int
	WaitingForPlayer string
}

func NewGame(gameId string) *Game {
	// Initlaises deck, playing field
	deck := NewDeck()
	pf := NewPlayingField(deck)

	// Constructs and returns the game
	return &Game{
		ID: gameId,

		Owner:   nil,
		Players: []*Player{},

		PlayingField: pf,
		Status:       InLobby,
		Round:        0,
		PlayOrder:    []int{},
	}

}

func (g *Game) findPlayerById(playerId string) *Player {
	for _, p := range g.Players {
		if p.ID == playerId {
			return p
		}
	}
	return nil
}

func (g *Game) AddPlayer(playerId string, playerConn *websocket.Conn) error {
	// Initilaises player with hands
	if g.Status != InLobby || g.Round != 0 {
		return &GameNotInLobbyError{g.ID}
	}

	// Constructing player
	player := NewPlayer(playerId, playerConn, g.PlayingField.Deck)

	// If no players in game, setting current player to owner
	if len(g.Players) == 0 {
		g.Owner = player
	}

	// Adding player to game
	g.Players = append(g.Players, player)

	return nil
}

func (g *Game) Start() error {
	if g.Status != InLobby {
		return &GameNotInLobbyError{gameId: g.ID}
	}

	g.Status = InGame

	// TODO: set playing order

	return nil
}

// Increments the game g
func (g *Game) Increment(player *Player, cmd PlayerCommand) error {
	if g.Status != InGame {
		return &GameNotStartedError{gameId: g.ID}
	}

	for _, card := range cmd.PlayCards {
		// Test to see if card can be played
		err := g.PlayingField.TestToPlayCard(card)
		if err != nil {
			return &CardCannotBePlayedError{
				CardToPlay:         card,
				CardOnPlayingField: g.PlayingField.Top(),
			}
		}

		// Playing card
		err = player.Play(card, g.PlayingField, g.Round)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) Run() error {
	// todo

	return nil
}

func (g *Game) Terminate() {
	for _, player := range g.Players {
		player.conn.Close()
	}
}

func (g *Game) handlePlayerMessage(player *Player, message []byte) error {
	// Tries to decode message as a game command
	var gameCommand GameCommand
	err := json.Unmarshal(message, &gameCommand)
	if err == nil {
		// Executes game command if player is game owner
		if player.ID != gameCommand.PlayerId {
			return &NotGameOwnerError{playerId: player.ID}
		} else if gameCommand.Action == "start" {
			g.Start()
		} else if gameCommand.Action == "terminate" {
			g.Terminate()
		} else {
			return &InvalidArgumentError{arg: gameCommand.Action}
		}
	}

	// Tries to decode message as player command
	var playerCommand PlayerCommand
	err = json.Unmarshal(message, &playerCommand)
	if err != nil {
		return &InvalidArgumentError{arg: message}
	}

	// Checks if it is player's turn
	if player.ID != g.WaitingForPlayer {
		return &NotPlayersTurnError{playerId: player.ID}
	}

	err = g.Increment(player, playerCommand)
	if err != nil {
		return err
	}

	return nil
}
