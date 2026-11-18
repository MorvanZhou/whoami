package gormdb

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"whoami/backend-go/internal/domain"
)

// StoreTokenRepo implements domain.StoreTokenRepository.
type StoreTokenRepo struct {
	db *gorm.DB
}

func NewStoreTokenRepo(db *gorm.DB) *StoreTokenRepo {
	return &StoreTokenRepo{db: db}
}

func (r *StoreTokenRepo) Create(ctx context.Context, token *domain.StoreApiToken) error {
	m := fromDomainStoreToken(token)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*token = *toDomainStoreToken(m)
	return nil
}

func (r *StoreTokenRepo) FindUnusedByToken(ctx context.Context, token string) (*domain.StoreApiToken, error) {
	var m storeTokenModel
	err := r.db.WithContext(ctx).
		Where("token = ? AND is_used = ?", token, false).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainStoreToken(m), nil
}

func (r *StoreTokenRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&storeTokenModel{}, "id = ?", id).Error
}

func (r *StoreTokenRepo) ExistsByAPIKeyID(ctx context.Context, apiKeyID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&storeTokenModel{}).
		Where("api_key_id = ? AND user_id = ?", apiKeyID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *StoreTokenRepo) CleanupExpired(ctx context.Context, cutoff time.Time) (int64, error) {
	res := r.db.WithContext(ctx).
		Where("is_used = ? OR created_at < ?", true, cutoff).
		Delete(&storeTokenModel{})
	return res.RowsAffected, res.Error
}
