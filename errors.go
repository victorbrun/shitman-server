package main

type CardNotInCollectionError struct{}

func (e *CardNotInCollectionError) Error() string {
	return "card not found in collection"
}

type CardAlreadyPlayedError struct{}

func (e *CardAlreadyPlayedError) Error() string {
	return "card has already been played"
}

type PEMDecodingError struct{}

func (e *PEMDecodingError) Error() string {
	return "failed to decode PEM block containing public key"
}

type NotRSAPublicKeyError struct{}

func (e *NotRSAPublicKeyError) Error() string {
	return "not an RSA public key"
}
