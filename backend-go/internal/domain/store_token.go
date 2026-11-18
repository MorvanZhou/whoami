package domain

import "time"

// StoreApiToken is a one-time token used to deliver a plaintext API key to an agent.
type StoreApiToken struct {
	ID          string
	UserID      string
	Token       string
	APIKeyPlain string
	APIKeyID    *string
	IsUsed      bool
	CreatedAt   time.Time
}
