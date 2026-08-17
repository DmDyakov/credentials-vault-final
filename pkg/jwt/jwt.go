// Package jwt предоставляет функции для работы с JWT-токенами.
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "credentials-vault"

type Token string

func (t Token) String() string {
	return string(t)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret    []byte
	accessTTL time.Duration
}

func New(secret string, accessTTL time.Duration) *Manager {
	return &Manager{
		secret:    []byte(secret),
		accessTTL: accessTTL,
	}
}

func (m *Manager) Generate(userID string) (Token, time.Time, error) {
	expiresAt := time.Now().Add(m.accessTTL)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return Token(signed), expiresAt, nil
}

// Verify проверяет и разбирает токен
func (m *Manager) Verify(token Token) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token.String(), &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.Issuer != issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	return claims, nil
}
