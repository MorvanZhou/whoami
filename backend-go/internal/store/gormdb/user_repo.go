package gormdb

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"whoami/backend-go/internal/domain"
)

// UserRepo implements domain.UserRepository.
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var m userModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(m), nil
}

func (r *UserRepo) FindByProvider(ctx context.Context, provider, providerID string) (*domain.User, error) {
	var m userModel
	err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_id = ?", provider, providerID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(m), nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m userModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(m), nil
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	m := fromDomainUser(user)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*user = *toDomainUser(m)
	return nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	m := fromDomainUser(user)
	return r.db.WithContext(ctx).Save(&m).Error
}
