package data

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenPayload struct {
	ID     string
	UserID string
	Email  string
}

type TokenDetails struct {
	Token     string
	ID        string
	Payload   TokenPayload
	ExpiresAt int64
}

type Claims struct {
	TokenID string `json:"token_id"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type RSACredentials struct {
	PrivateKey []byte
	PublicKey  []byte
	TTL        time.Duration
}
