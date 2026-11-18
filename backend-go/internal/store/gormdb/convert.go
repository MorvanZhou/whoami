package gormdb

import "whoami/backend-go/internal/domain"

func toDomainUser(m userModel) *domain.User {
	return &domain.User{
		ID:         m.ID,
		Email:      m.Email,
		Name:       m.Name,
		AvatarURL:  m.AvatarURL,
		Provider:   m.Provider,
		ProviderID: m.ProviderID,
		CreatedAt:  m.CreatedAt,
		IsActive:   m.IsActive,
	}
}

func fromDomainUser(u *domain.User) userModel {
	return userModel{
		ID:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		AvatarURL:  u.AvatarURL,
		Provider:   u.Provider,
		ProviderID: u.ProviderID,
		CreatedAt:  u.CreatedAt,
		IsActive:   u.IsActive,
	}
}

func toDomainAPIKey(m apiKeyModel) *domain.ApiKey {
	return &domain.ApiKey{
		ID:        m.ID,
		UserID:    m.UserID,
		KeyHash:   m.KeyHash,
		KeyPrefix: m.KeyPrefix,
		KeySuffix: m.KeySuffix,
		Label:     m.Label,
		CreatedAt: m.CreatedAt,
		LastUsed:  m.LastUsed,
		IsActive:  m.IsActive,
	}
}

func fromDomainAPIKey(k *domain.ApiKey) apiKeyModel {
	return apiKeyModel{
		ID:        k.ID,
		UserID:    k.UserID,
		KeyHash:   k.KeyHash,
		KeyPrefix: k.KeyPrefix,
		KeySuffix: k.KeySuffix,
		Label:     k.Label,
		CreatedAt: k.CreatedAt,
		LastUsed:  k.LastUsed,
		IsActive:  k.IsActive,
	}
}

func toDomainStoreToken(m storeTokenModel) *domain.StoreApiToken {
	return &domain.StoreApiToken{
		ID:          m.ID,
		UserID:      m.UserID,
		Token:       m.Token,
		APIKeyPlain: m.APIKeyPlain,
		APIKeyID:    m.APIKeyID,
		IsUsed:      m.IsUsed,
		CreatedAt:   m.CreatedAt,
	}
}

func fromDomainStoreToken(t *domain.StoreApiToken) storeTokenModel {
	return storeTokenModel{
		ID:          t.ID,
		UserID:      t.UserID,
		Token:       t.Token,
		APIKeyPlain: t.APIKeyPlain,
		APIKeyID:    t.APIKeyID,
		IsUsed:      t.IsUsed,
		CreatedAt:   t.CreatedAt,
	}
}
