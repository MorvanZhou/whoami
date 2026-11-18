package domain

import "time"

// User is a whoami account linked to an OAuth identity.
type User struct {
	ID         string
	Email      *string
	Name       *string
	AvatarURL  *string
	Provider   string
	ProviderID string
	CreatedAt  time.Time
	IsActive   bool
}
