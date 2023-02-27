package token

import "time"

type Maker interface {

	// CreateToken creates a new token for a specify username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyTOken check if token is valid or not
	VerifyToken(token string) (*Payload, error)
}
