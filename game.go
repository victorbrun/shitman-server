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
	Card          Card  `json:"card"`
	PlayedAt      int64 // Unix timestamp for when card was played
	PlayedBy      *Player
	PlayedInRound int
}

func NewPlayer(id string, deck *Deck) *Player {
	return &Player{
		ID:          id,
		PublicHand:  NewHand(3, deck),
		PrivateHand: NewHand(3, deck),
		HiddenHand:  NewHand(3, deck),
	}
}

func NewPlayedCard(card Card, playedBy *Player, playedInRound int) PlayedCard {
	return PlayedCard{
		Card:          card,
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
