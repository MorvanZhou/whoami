package gormdb

import "time"

type userModel struct {
	ID         string    `gorm:"column:id;primaryKey"`
	Email      *string   `gorm:"column:email;uniqueIndex"`
	Name       *string   `gorm:"column:name"`
	AvatarURL  *string   `gorm:"column:avatar_url"`
	Provider   string    `gorm:"column:provider;not null;uniqueIndex:uq_provider_identity"`
	ProviderID string    `gorm:"column:provider_id;not null;uniqueIndex:uq_provider_identity"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	IsActive   bool      `gorm:"column:is_active;default:true"`
}

func (userModel) TableName() string { return "users" }

type apiKeyModel struct {
	ID        string     `gorm:"column:id;primaryKey"`
	UserID    string     `gorm:"column:user_id;not null;index"`
	KeyHash   string     `gorm:"column:key_hash;uniqueIndex;not null"`
	KeyPrefix string     `gorm:"column:key_prefix;not null"`
	KeySuffix string     `gorm:"column:key_suffix;not null;default:''"`
	Label     *string    `gorm:"column:label"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	LastUsed  *time.Time `gorm:"column:last_used"`
	IsActive  bool       `gorm:"column:is_active;default:true"`
}

func (apiKeyModel) TableName() string { return "api_keys" }

type storeTokenModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	UserID      string    `gorm:"column:user_id;not null;index"`
	Token       string    `gorm:"column:token;uniqueIndex;not null"`
	APIKeyPlain string    `gorm:"column:api_key_plain;not null"`
	APIKeyID    *string   `gorm:"column:api_key_id"`
	IsUsed      bool      `gorm:"column:is_used;default:false"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (storeTokenModel) TableName() string { return "store_api_tokens" }
