package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/romankravchuk/eldorado/internal/data"
)

func CreateToken(payload *data.TokenPayload, ttl time.Duration, prvKey []byte) (*data.TokenDetails, error) {
	now := time.Now().UTC()
	td := &data.TokenDetails{
		ID:        uuid.New().String(),
		Payload:   *payload,
		ExpiresAt: now.Add(ttl).Unix(),
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return nil, err
	}
	claims := data.Claims{
		TokenID: td.ID,
		UserID:  payload.ID,
		Email:   payload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        td.ID,
		},
	}

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, &claims).SignedString(key)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ValidateToken(token string, publickKey []byte) (*data.TokenPayload, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publickKey)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		&data.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*data.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	payload := &data.TokenPayload{
		ID:     claims.TokenID,
		UserID: claims.UserID,
		Email:  claims.Email,
	}

	return payload, nil
}
