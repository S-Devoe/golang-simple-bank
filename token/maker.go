package token

import "time"

// maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, email string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
