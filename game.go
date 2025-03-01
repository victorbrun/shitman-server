package main

import (
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

type Game struct {
	// Metadata
	ID string

	// Player data
	Owner   *Player
	Players []*Player

	// Game data
	PlayingField *PlayingField
	Status       GameStatus
	Round        int
	PlayOrder    []int
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
func (g *Game) Increment(playedCards []PlayedCard) error {
	if g.Status != InGame {
		return &GameNotStartedError{gameId: g.ID}
	}

	// TODO: play cards

	return nil
}

func (g *Game) Run() error {
	// todo

	return nil
}
