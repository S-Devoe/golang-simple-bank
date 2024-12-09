package token

import (
	"errors"
	"time"

	"github.com/S-Devoe/golang-simple-bank/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// payload contains  the payload data of the token
type Payload struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	ID        ulid.ULID `json:"id"`
	Audience  string    `json:"audience,omitempty"`
}

func NewPayload(username, email string, duration time.Duration) (*Payload, error) {
	id, err := util.GenerateULID()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Username:  username,
		Email:     email,
		ExpiredAt: time.Now().Add(duration),
		IssuedAt:  time.Now(),
		ID:        id,
	}

	return payload, nil
}

// check if token is valid
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	if p.Audience == "" {
		return nil, nil // No audience set
	}
	return jwt.ClaimStrings{p.Audience}, nil
}

// Implement additional methods as needed (optional)
func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetExpiresAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetSubject returns the subject claim
func (p *Payload) GetSubject() (string, error) {
	return "", nil
}
