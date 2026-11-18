package service_test

import (
	"context"
	"sync"
	"time"

	"whoami/backend-go/internal/domain"
)

type mockProfileStore struct {
	mu       sync.Mutex
	profiles map[string]string
	history  map[string][]domain.HistoryItem
	contents map[string]map[string]string // userID -> ts -> content

	getErr      error
	saveErr     error
	listErr     error
	rollbackErr error
}

func newMockProfileStore() *mockProfileStore {
	return &mockProfileStore{
		profiles: make(map[string]string),
		history:  make(map[string][]domain.HistoryItem),
		contents: make(map[string]map[string]string),
	}
}

func (m *mockProfileStore) Get(_ context.Context, userID string) (string, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.getErr != nil {
		return "", false, m.getErr
	}
	c, ok := m.profiles[userID]
	return c, ok, nil
}

func (m *mockProfileStore) Save(_ context.Context, userID, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.saveErr != nil {
		return m.saveErr
	}
	if old, ok := m.profiles[userID]; ok {
		ts := time.Now().UTC().Format("20060102T150405Z")
		m.history[userID] = append(m.history[userID], domain.HistoryItem{
			Timestamp:    ts,
			Size:         int64(len(old)),
			LastModified: time.Now().UTC().Format(time.RFC3339),
		})
		if m.contents[userID] == nil {
			m.contents[userID] = make(map[string]string)
		}
		m.contents[userID][ts] = old
	}
	m.profiles[userID] = content
	return nil
}

func (m *mockProfileStore) ListHistory(_ context.Context, userID string) ([]domain.HistoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listErr != nil {
		return nil, m.listErr
	}
	items := append([]domain.HistoryItem(nil), m.history[userID]...)
	// newest first
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return items, nil
}

func (m *mockProfileStore) Rollback(_ context.Context, userID, timestamp string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.rollbackErr != nil {
		return false, m.rollbackErr
	}
	content, ok := m.contents[userID][timestamp]
	if !ok {
		return false, nil
	}
	m.profiles[userID] = content
	return true, nil
}

type mockUserRepo struct {
	mu    sync.Mutex
	byID  map[string]*domain.User
	err   error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{byID: make(map[string]*domain.User)}
}

func (m *mockUserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return nil, m.err
	}
	u, ok := m.byID[id]
	if !ok {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) FindByProvider(context.Context, string, string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByEmail(context.Context, string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Create(_ context.Context, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *user
	m.byID[user.ID] = &cp
	return nil
}
func (m *mockUserRepo) Update(_ context.Context, user *domain.User) error {
	return m.Create(context.Background(), user)
}

type mockAPIKeyRepo struct {
	mu   sync.Mutex
	keys map[string]*domain.ApiKey // id -> key
	err  error
}

func newMockAPIKeyRepo() *mockAPIKeyRepo {
	return &mockAPIKeyRepo{keys: make(map[string]*domain.ApiKey)}
}

func (m *mockAPIKeyRepo) Create(_ context.Context, key *domain.ApiKey) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	cp := *key
	m.keys[key.ID] = &cp
	return nil
}

func (m *mockAPIKeyRepo) ListActiveByUser(_ context.Context, userID string) ([]domain.ApiKey, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []domain.ApiKey
	for _, k := range m.keys {
		if k.UserID == userID && k.IsActive {
			out = append(out, *k)
		}
	}
	return out, nil
}

func (m *mockAPIKeyRepo) FindActiveByHash(_ context.Context, keyHash string) (*domain.ApiKey, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range m.keys {
		if k.KeyHash == keyHash && k.IsActive {
			cp := *k
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *mockAPIKeyRepo) DeleteByIDForUser(_ context.Context, keyID, userID string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	k, ok := m.keys[keyID]
	if !ok || k.UserID != userID {
		return false, nil
	}
	delete(m.keys, keyID)
	return true, nil
}

func (m *mockAPIKeyRepo) TouchLastUsed(_ context.Context, keyID string, at time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if k, ok := m.keys[keyID]; ok {
		k.LastUsed = &at
	}
	return nil
}

type mockStoreTokenRepo struct {
	mu     sync.Mutex
	tokens map[string]*domain.StoreApiToken // id -> token
	err    error
}

func newMockStoreTokenRepo() *mockStoreTokenRepo {
	return &mockStoreTokenRepo{tokens: make(map[string]*domain.StoreApiToken)}
}

func (m *mockStoreTokenRepo) Create(_ context.Context, token *domain.StoreApiToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := *token
	m.tokens[token.ID] = &cp
	return nil
}

func (m *mockStoreTokenRepo) FindUnusedByToken(_ context.Context, token string) (*domain.StoreApiToken, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return nil, m.err
	}
	for _, t := range m.tokens {
		if t.Token == token && !t.IsUsed {
			cp := *t
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *mockStoreTokenRepo) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, id)
	return nil
}

func (m *mockStoreTokenRepo) ExistsByAPIKeyID(_ context.Context, apiKeyID, userID string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range m.tokens {
		if t.UserID == userID && t.APIKeyID != nil && *t.APIKeyID == apiKeyID {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockStoreTokenRepo) CleanupExpired(_ context.Context, cutoff time.Time) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var n int64
	for id, t := range m.tokens {
		if t.IsUsed || t.CreatedAt.Before(cutoff) {
			delete(m.tokens, id)
			n++
		}
	}
	return n, nil
}
