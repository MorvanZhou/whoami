package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const algorithm = "HS256"

// TokenManager issues and validates HS256 JWTs (sub = user id).
type TokenManager struct {
	secret     []byte
	expireDays int
}

func NewTokenManager(secret string, expireDays int) *TokenManager {
	if expireDays <= 0 {
		expireDays = 7
	}
	return &TokenManager{
		secret:     []byte(secret),
		expireDays: expireDays,
	}
}

func (m *TokenManager) Create(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().UTC().Add(time.Duration(m.expireDays) * 24 * time.Hour).Unix(),
		"iat": time.Now().UTC().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *TokenManager) ParseUserID(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != algorithm {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return "", fmt.Errorf("missing subject")
	}
	return sub, nil
}

func (m *TokenManager) MaxAgeSeconds() int {
	return m.expireDays * 86400
}
