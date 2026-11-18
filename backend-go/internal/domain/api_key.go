package domain

import "time"

// ApiKey is a hashed agent credential belonging to a user.
type ApiKey struct {
	ID        string
	UserID    string
	KeyHash   string
	KeyPrefix string
	KeySuffix string
	Label     *string
	CreatedAt time.Time
	LastUsed  *time.Time
	IsActive  bool
}
