package service

import (
	"context"
	"fmt"
	"log/slog"
	"unicode/utf8"

	"whoami/backend-go/internal/domain"
)

const ProfileMaxLength = 5000

// ProfileService exposes profile read/write over a ProfileStore.
type ProfileService struct {
	store domain.ProfileStore
}

func NewProfileService(store domain.ProfileStore) *ProfileService {
	return &ProfileService{store: store}
}

func (s *ProfileService) Get(ctx context.Context, userID string) (*string, error) {
	content, ok, err := s.store.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		slog.Info("get profile", "user_id", userID, "found", false)
		return nil, nil
	}
	slog.Info("get profile", "user_id", userID, "found", true)
	return &content, nil
}

func (s *ProfileService) Save(ctx context.Context, userID, content string) error {
	if utf8.RuneCountInString(content) > ProfileMaxLength {
		return fmt.Errorf("content exceeds %d characters", ProfileMaxLength)
	}
	if err := s.store.Save(ctx, userID, content); err != nil {
		return err
	}
	slog.Info("update profile", "user_id", userID, "len", utf8.RuneCountInString(content))
	return nil
}

func (s *ProfileService) ListHistory(ctx context.Context, userID string) ([]domain.HistoryItem, error) {
	return s.store.ListHistory(ctx, userID)
}

func (s *ProfileService) Rollback(ctx context.Context, userID, timestamp string) (bool, error) {
	return s.store.Rollback(ctx, userID, timestamp)
}
