package domain

import (
	"context"
	"time"
)

// UserRepository persists users.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByProvider(ctx context.Context, provider, providerID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

// ApiKeyRepository persists API keys.
type ApiKeyRepository interface {
	Create(ctx context.Context, key *ApiKey) error
	ListActiveByUser(ctx context.Context, userID string) ([]ApiKey, error)
	FindActiveByHash(ctx context.Context, keyHash string) (*ApiKey, error)
	DeleteByIDForUser(ctx context.Context, keyID, userID string) (bool, error)
	TouchLastUsed(ctx context.Context, keyID string, at time.Time) error
}

// StoreTokenRepository persists one-time store tokens.
type StoreTokenRepository interface {
	Create(ctx context.Context, token *StoreApiToken) error
	FindUnusedByToken(ctx context.Context, token string) (*StoreApiToken, error)
	Delete(ctx context.Context, id string) error
	ExistsByAPIKeyID(ctx context.Context, apiKeyID, userID string) (bool, error)
	CleanupExpired(ctx context.Context, cutoff time.Time) (int64, error)
}
