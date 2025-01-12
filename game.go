package main

type GameStatus int

const (
	InLobby GameStatus = iota
	InGame
	Finished
	Canceled
)

type Player struct {
	// Metadata
	ID string

	// Cards
	PrivateHand *Hand
	PublicHand  *Hand
	HiddenHand  *Hand
}

func NewPlayer(id string, deck *Deck) *Player {
	return &Player{
		ID:          id,
		PublicHand:  NewHand(3, deck),
		PrivateHand: NewHand(3, deck),
		HiddenHand:  NewHand(3, deck),
	}
}

// Consors the hand such that the information can be
// communicated to every participant in the game
func (p *Player) PublicCensor() *Player {
	return &Player{
		ID:          p.ID,
		PrivateHand: nil,
		PublicHand:  p.PublicHand,
		HiddenHand:  nil,
	}
}

// Censors the hands such that the information can be communicated
// to the player owning this hand
func (p *Player) PrivateCensor() *Player {
	return &Player{
		ID:          p.ID,
		PrivateHand: p.PrivateHand,
		PublicHand:  p.PublicHand,
		HiddenHand:  nil,
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

// Record of the state of a game at the end of a round
type GameState struct {
	Players      []*Player
	PlayingField *PlayingField
	Round        int
}

type Game struct {
	ID         string
	GameStates []*GameState
	Status     GameState
}

func NewGame(playerIDs []string) *Game {
	// Initlaises deck, playing field
	deck := NewDeck()
	pf := NewPlayingField(deck)

	// Initilaises players with their respective hands
	players := make([]*Player, len(playerIDs))
	for ix, playerID := range playerIDs {
		players[ix] = NewPlayer(playerID, deck)
	}

	// Manually creates the initial game state
	initialGameState := &GameState{
		Players:      players,
		PlayingField: pf,
		Round:        0,
	}

	// Constructs and returns the game
	var gameStates []*GameState = []*GameState{initialGameState}
	return &Game{
		GameStates: gameStates,
	}

}

// Increments the game g
func (g *Game) IncrementGame(playedCards []PlayedCard) {

}
