package main

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

type PlayerAlreadyConnectedError struct{}

func (e *PlayerAlreadyConnectedError) Error() string {
	return "player is already connected"
}
