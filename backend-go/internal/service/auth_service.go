package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"whoami/backend-go/internal/domain"
)

// AuthService handles user identity lookup/creation.
type AuthService struct {
	users domain.UserRepository
}

func NewAuthService(users domain.UserRepository) *AuthService {
	return &AuthService{users: users}
}

// FindOrCreateUser mirrors the Python auth_service.find_or_create_user behavior.
func (s *AuthService) FindOrCreateUser(
	ctx context.Context,
	provider, providerID string,
	email, name, avatarURL *string,
) (*domain.User, error) {
	user, err := s.users.FindByProvider(ctx, provider, providerID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		changed := false
		if email != nil {
			user.Email = email
			changed = true
		}
		if name != nil {
			user.Name = name
			changed = true
		}
		if avatarURL != nil {
			user.AvatarURL = avatarURL
			changed = true
		}
		if changed {
			if err := s.users.Update(ctx, user); err != nil {
				return nil, err
			}
		}
		slog.Info("existing user logged in", "user_id", user.ID, "provider", provider)
		return user, nil
	}

	// Same email from a different provider — link to existing account (Python behavior).
	if email != nil && *email != "" {
		existing, err := s.users.FindByEmail(ctx, *email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			if name != nil && (existing.Name == nil || *existing.Name == "") {
				existing.Name = name
			}
			if avatarURL != nil && (existing.AvatarURL == nil || *existing.AvatarURL == "") {
				existing.AvatarURL = avatarURL
			}
			if err := s.users.Update(ctx, existing); err != nil {
				return nil, err
			}
			slog.Info("user linked by email",
				"user_id", existing.ID,
				"provider", provider,
				"original_provider", existing.Provider,
			)
			return existing, nil
		}
	}

	now := time.Now().UTC()
	user = &domain.User{
		ID:         uuid.NewString(),
		Email:      email,
		Name:       name,
		AvatarURL:  avatarURL,
		Provider:   provider,
		ProviderID: providerID,
		CreatedAt:  now,
		IsActive:   true,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	slog.Info("new user created", "user_id", user.ID, "provider", provider)
	return user, nil
}
