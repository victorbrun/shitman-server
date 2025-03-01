package main

import "fmt"

type CardNotInCollectionError struct{}

func (e *CardNotInCollectionError) Error() string {
	return "card not found in collection"
}

type CardAlreadyPlayedError struct{}

func (e *CardAlreadyPlayedError) Error() string {
	return "card has already been played"
}

type GameNotInMapError struct{}

func (e *GameNotInMapError) Error() string {
	return "game is not in map"
}

type GameNotInLobbyError struct{ gameId string }

func (e *GameNotInLobbyError) Error() string {
	return fmt.Sprintf("game (%v) not in lobby", e.gameId)
}

type GameNotStartedError struct{ gameId string }

func (e *GameNotStartedError) Error() string {
	return fmt.Sprintf("game (%v) not started", e.gameId)
}

type PlayerAlreadyConnectedError struct{}

func (e *PlayerAlreadyConnectedError) Error() string {
	return "player is already connected"
}

type InvalidArgumentError struct {
	arg any
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("following argument is not valid: %v", e.arg)
}
