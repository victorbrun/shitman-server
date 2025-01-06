package main

import "crypto/rsa"

type Player struct {
	// Metadata
	ID        string
	PublicKey *rsa.PublicKey

	// Cards
	Hand           *Hand
	DownFacingHand *Hand
	UpFacingHand   *Hand
}

func NewPlayer(id string, deck *Deck) *Player {
	return &Player{
		ID:             id,
		Hand:           NewHand(3, deck),
		UpFacingHand:   NewHand(3, deck),
		DownFacingHand: NewHand(3, deck),
	}
}

type Game struct {
	Players            []*Player
	ActivePlayingField *PlayingField
	TurnedCards        *PlayingField
	Deck               *Deck
	CurrentRound       int
}

func NewGame(playerIDs []string) *Game {
	// Initlaises deck, playing field
	deck := NewDeck()
	pf := NewPlayingField()

	// Initilaises players with their respective hands
	players := make([]*Player, len(playerIDs))
	for ix, playerID := range playerIDs {
		players[ix] = NewPlayer(playerID, deck)
	}

	return &Game{
		Players:            players,
		ActivePlayingField: pf,
		Deck:               deck,
		CurrentRound:       0,
	}
}

type PublicGameState struct {
	Players []struct {
		ID                    string `json:"id"`
		NumberCardsOnHand     int    `json:"number_cars_on_hand"`
		NumberDownFacingCards int    `json:"number_down_facing_cards"`
		UpFacingHand          []Card `json:"up_facing_cards"`
	}
	PlayedCards []Card `json:"played_cards"`
}
