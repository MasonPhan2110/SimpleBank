package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	publicKey    ed25519.PublicKey
	privateKey   ed25519.PrivateKey
}

// Create new PasetoMaker
func NewPasetoMaker(symmertricKey string) (Maker, error) {
	if len(symmertricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must have exactly %d characters", chacha20poly1305.KeySize)
	}
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmertricKey),
		publicKey:    publicKey,
		privateKey:   privateKey,
	}

	return maker, nil
}

// CreateToken creates a new token for a specify username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyTOken check if token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (maker *PasetoMaker) CreateTokenPublic(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	footer := "footer"

	token, err := maker.paseto.Sign(maker.privateKey, payload, footer)
	if err != nil {
		return "", err
	}
	return token, nil
}
