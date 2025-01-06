package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
)

type PublicPlayerData struct {
	ID                    string `json:"id"`
	NumberCardsOnHand     int    `json:"number_cars_on_hand"`
	NumberDownFacingCards int    `json:"number_down_facing_cards"`
	EncryptedHand         []byte `json:"encrypted_hand"`
	UpFacingHand          []Card `json:"up_facing_cards"`
}

// Creates a public player data struct from player
func NewPublicPlayerData(p *Player) *PublicPlayerData {
	// Extracts and converts the hand to json
	privateHandBytes, err := json.Marshal(p.Hand.Cards)
	if err != nil {
		panic(err)
	}

	// Encrypts them using players public key
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, p.PublicKey, privateHandBytes, nil)

	// Constructs and returns the struct
	return &PublicPlayerData{
		ID:                    p.ID,
		NumberCardsOnHand:     p.Hand.Size(),
		NumberDownFacingCards: p.DownFacingHand.Size(),
		EncryptedHand:         ciphertext,
		UpFacingHand:          p.UpFacingHand.Cards,
	}
}

func parsePublicKeyPEM(publicKeyPEMBites []byte) (*rsa.PublicKey, error) {
	// Decode the PEM string
	block, _ := pem.Decode(publicKeyPEMBites)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, &PEMDecodingError{}
	}

	// Parse the public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, &NotRSAPublicKeyError{}
	}

	return rsaPubKey, nil
}
