package gormdb

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"whoami/backend-go/internal/domain"
)

// ApiKeyRepo implements domain.ApiKeyRepository.
type ApiKeyRepo struct {
	db *gorm.DB
}

func NewApiKeyRepo(db *gorm.DB) *ApiKeyRepo {
	return &ApiKeyRepo{db: db}
}

func (r *ApiKeyRepo) Create(ctx context.Context, key *domain.ApiKey) error {
	m := fromDomainAPIKey(key)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*key = *toDomainAPIKey(m)
	return nil
}

func (r *ApiKeyRepo) ListActiveByUser(ctx context.Context, userID string) ([]domain.ApiKey, error) {
	var models []apiKeyModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.ApiKey, 0, len(models))
	for _, m := range models {
		out = append(out, *toDomainAPIKey(m))
	}
	return out, nil
}

func (r *ApiKeyRepo) FindActiveByHash(ctx context.Context, keyHash string) (*domain.ApiKey, error) {
	var m apiKeyModel
	err := r.db.WithContext(ctx).
		Where("key_hash = ? AND is_active = ?", keyHash, true).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainAPIKey(m), nil
}

func (r *ApiKeyRepo) DeleteByIDForUser(ctx context.Context, keyID, userID string) (bool, error) {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", keyID, userID).
		Delete(&apiKeyModel{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *ApiKeyRepo) TouchLastUsed(ctx context.Context, keyID string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&apiKeyModel{}).
		Where("id = ?", keyID).
		Update("last_used", at).Error
}
