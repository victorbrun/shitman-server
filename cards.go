package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type Suit string
type Rank string

const (
	Clubs    Suit = "Clubs"
	Diamonds Suit = "Diamonds"
	Hearts   Suit = "Hearts"
	Spades   Suit = "Spades"

	Ace   Rank = "Ace"
	Two   Rank = "Two"
	Three Rank = "Three"
	Four  Rank = "Four"
	Five  Rank = "Five"
	Six   Rank = "Six"
	Seven Rank = "Seven"
	Eight Rank = "Eight"
	Nine  Rank = "Nine"
	Ten   Rank = "Ten"
	Jack  Rank = "Jack"
	Queen Rank = "Queen"
	King  Rank = "King"
)

func (s *Suit) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch {
	case str == "Clubs":
		*s = Clubs
	case str == "Diamonds":
		*s = Diamonds
	case str == "Hearts":
		*s = Hearts
	case str == "Spades":
		*s = Spades
	default:
		return &InvalidArgumentError{arg: str}
	}
	return nil
}

func (r *Rank) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch {
	case str == "Ace" || str == "1":
		*r = Ace
	case str == "Two" || str == "2":
		*r = Two
	case str == "Tree" || str == "3":
		*r = Three
	case str == "Four" || str == "4":
		*r = Four
	case str == "Five" || str == "5":
		*r = Five
	case str == "Six" || str == "6":
		*r = Six
	case str == "Seven" || str == "7":
		*r = Seven
	case str == "Eight" || str == "8":
		*r = Eight
	case str == "Nine" || str == "9":
		*r = Nine
	case str == "Ten" || str == "10":
		*r = Ten
	case str == "Jack" || str == "11":
		*r = Jack
	case str == "Queen" || str == "12":
		*r = Queen
	case str == "King" || str == "13":
		*r = King
	default:
		return &InvalidArgumentError{arg: str}
	}
	return nil
}

// Returns the numeric value of the
// rank such that cards can be compared
func (r Rank) Value(aceIsMostValued bool) int {
	switch r {
	case Ace:
		if aceIsMostValued {
			return 14
		}
		return 1

	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten:
		return 10
	case Jack:
		return 11
	case Queen:
		return 12
	case King:
		return 13
	default:
		panic("something went wrong this code should not be reachable")
	}
}

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

func (c Card) String() string {
	return fmt.Sprintf("{%v of %v}", c.Rank, c.Suit)
}

type Collection struct {
	Cards []Card `json:"cards"`
}

// Returns the number of cards in the collection
func (c Collection) Size() int {
	return len(c.Cards)
}

// Adds collection of cards newCards to collection
// of cards c
func (c *Collection) Merge(newCards Collection) {
	c.Cards = append(c.Cards, newCards.Cards...)
}

// Removes the card at index idx from c
func (c *Collection) Remove(idx int) {
	c.Cards[idx] = c.Cards[c.Size()-1]
	c.Cards = c.Cards[:c.Size()-1]
}

// Checks if card is contained within c.
// If it is the function returns (index of card in c.Cards, true)
// If it is not the functions returns (-1, false)
func (c *Collection) Contains(card Card) (int, bool) {
	for ix, cardInCollection := range c.Cards {
		if card == cardInCollection {
			return ix, true
		}
	}
	return -1, false
}

type Deck struct {
	Collection
}

func NewDeck() *Deck {
	// The deck consists of the cartetian product
	// between all suits and ranks put in a flat matrix
	suits := []Suit{Clubs, Diamonds, Hearts, Spades}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
	cards := make([]Card, 52)
	for ix, suit := range suits {
		for jx, rank := range ranks {
			cards[13*ix+jx] = Card{Suit: suit, Rank: rank}
		}
	}

	cardCollection := Collection{Cards: cards}
	return &Deck{Collection: cardCollection}
}

// Draws n random cards from the deck and removes
// the drawn cards from the deck
func (d *Deck) Draw(n int) Collection {
	drawnCards := make([]Card, n)
	for ix := 0; ix < n; ix++ {
		// Randomly selects a card from the deck
		randomIndex := rand.IntN(d.Size())
		drawnCards[ix] = d.Cards[randomIndex]

		// Removes the selected card from the deck
		d.Remove(randomIndex)
	}

	return Collection{Cards: drawnCards}
}

type Hand struct {
	Collection
}

// Creates a new hand by drawing n cards from deck
func NewHand(n int, deck *Deck) *Hand {
	return &Hand{Collection: deck.Draw(n)}
}

// Draws n new cards from deck and adds them to h
func (h *Hand) drawNewCards(n int, deck Deck) {
	newCards := deck.Draw(n)
	h.Merge(newCards)
}
