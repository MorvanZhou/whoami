package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"whoami/backend-go/internal/domain"
)

const StoreTokenExpireMinutes = 15

// StoreTokenService manages one-time API key delivery tokens.
type StoreTokenService struct {
	tokens domain.StoreTokenRepository
}

func NewStoreTokenService(tokens domain.StoreTokenRepository) *StoreTokenService {
	return &StoreTokenService{tokens: tokens}
}

func (s *StoreTokenService) Create(ctx context.Context, userID, apiKeyPlain, apiKeyID string) (string, error) {
	token, err := generateStoreToken()
	if err != nil {
		return "", err
	}
	id := apiKeyID
	rec := &domain.StoreApiToken{
		ID:          uuid.NewString(),
		UserID:      userID,
		Token:       token,
		APIKeyPlain: apiKeyPlain,
		APIKeyID:    &id,
		IsUsed:      false,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.tokens.Create(ctx, rec); err != nil {
		return "", err
	}
	slog.Info("store token created", "user_id", userID)
	return token, nil
}

// Consume returns the plaintext API key and deletes the token. Returns ok=false if missing/expired.
func (s *StoreTokenService) Consume(ctx context.Context, token string) (plainKey string, ok bool, err error) {
	rec, err := s.tokens.FindUnusedByToken(ctx, token)
	if err != nil {
		return "", false, err
	}
	if rec == nil {
		return "", false, nil
	}

	createdAt := rec.CreatedAt
	if createdAt.Location() != time.UTC {
		createdAt = createdAt.UTC()
	}
	// SQLite may return naive times; treat as UTC.
	expireAt := createdAt.Add(StoreTokenExpireMinutes * time.Minute)
	if time.Now().UTC().After(expireAt) {
		_ = s.tokens.Delete(ctx, rec.ID)
		slog.Info("store token expired and deleted", "token_prefix", safePrefix(token))
		return "", false, nil
	}

	plain := rec.APIKeyPlain
	if err := s.tokens.Delete(ctx, rec.ID); err != nil {
		return "", false, err
	}
	slog.Info("store token consumed", "user_id", rec.UserID)
	return plain, true, nil
}

// IsConsumed reports whether the store token for an api key is gone (configured by agent).
func (s *StoreTokenService) IsConsumed(ctx context.Context, apiKeyID, userID string) (bool, error) {
	exists, err := s.tokens.ExistsByAPIKeyID(ctx, apiKeyID, userID)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

func (s *StoreTokenService) CleanupExpired(ctx context.Context) (int64, error) {
	cutoff := time.Now().UTC().Add(-StoreTokenExpireMinutes * time.Minute)
	return s.tokens.CleanupExpired(ctx, cutoff)
}

func generateStoreToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func safePrefix(token string) string {
	if len(token) <= 8 {
		return token
	}
	return token[:8]
}
