package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"whoami/backend-go/internal/domain"
)

// KeyService manages API keys.
type KeyService struct {
	keys         domain.ApiKeyRepository
	users        domain.UserRepository
	apiKeyPrefix string
}

func NewKeyService(keys domain.ApiKeyRepository, users domain.UserRepository, apiKeyPrefix string) *KeyService {
	if apiKeyPrefix == "" {
		apiKeyPrefix = "wai_"
	}
	return &KeyService{keys: keys, users: users, apiKeyPrefix: apiKeyPrefix}
}

func (s *KeyService) Create(ctx context.Context, userID, label string) (plainKey string, key *domain.ApiKey, err error) {
	plainKey, err = generateAPIKey(s.apiKeyPrefix)
	if err != nil {
		return "", nil, err
	}
	lbl := label
	key = &domain.ApiKey{
		ID:        uuid.NewString(),
		UserID:    userID,
		KeyHash:   hashAPIKey(plainKey),
		KeyPrefix: plainKey[:min(8, len(plainKey))],
		KeySuffix: plainKey[max(0, len(plainKey)-4):],
		Label:     &lbl,
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}
	if err := s.keys.Create(ctx, key); err != nil {
		return "", nil, err
	}
	slog.Info("api key created", "user_id", userID, "prefix", key.KeyPrefix)
	return plainKey, key, nil
}

func (s *KeyService) List(ctx context.Context, userID string) ([]domain.ApiKey, error) {
	return s.keys.ListActiveByUser(ctx, userID)
}

func (s *KeyService) Revoke(ctx context.Context, keyID, userID string) (bool, error) {
	ok, err := s.keys.DeleteByIDForUser(ctx, keyID, userID)
	if err != nil {
		return false, err
	}
	if ok {
		slog.Info("api key deleted", "key_id", keyID)
	}
	return ok, nil
}

// Authenticate validates a plaintext API key and returns the owning active user.
func (s *KeyService) Authenticate(ctx context.Context, plainKey string) (*domain.User, error) {
	apiKey, err := s.keys.FindActiveByHash(ctx, hashAPIKey(plainKey))
	if err != nil {
		return nil, err
	}
	if apiKey == nil {
		return nil, nil
	}
	_ = s.keys.TouchLastUsed(ctx, apiKey.ID, time.Now().UTC())

	user, err := s.users.FindByID(ctx, apiKey.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsActive {
		return nil, nil
	}
	return user, nil
}

func generateAPIKey(prefix string) (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate api key: %w", err)
	}
	return prefix + hex.EncodeToString(buf), nil
}

func hashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
