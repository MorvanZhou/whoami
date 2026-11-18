package auth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

// StateStore is an in-memory OAuth state cache with TTL (CSRF protection).
type StateStore struct {
	mu      sync.Mutex
	ttl     time.Duration
	maxSize int
	items   map[string]stateEntry
}

type stateEntry struct {
	value     string
	expiresAt time.Time
}

func NewStateStore(ttl time.Duration, maxSize int) *StateStore {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	if maxSize <= 0 {
		maxSize = 10000
	}
	return &StateStore{
		ttl:     ttl,
		maxSize: maxSize,
		items:   make(map[string]stateEntry),
	}
}

// Put stores value under a newly generated state key and returns the key.
func (s *StateStore) Put(value string) (string, error) {
	key, err := randomState()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpiredLocked(time.Now())
	if len(s.items) >= s.maxSize {
		// Drop an arbitrary entry to bound memory.
		for k := range s.items {
			delete(s.items, k)
			break
		}
	}
	s.items[key] = stateEntry{value: value, expiresAt: time.Now().Add(s.ttl)}
	return key, nil
}

// Take returns and deletes the value for state. ok=false if missing/expired.
func (s *StateStore) Take(state string) (value string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.purgeExpiredLocked(now)
	e, exists := s.items[state]
	if !exists {
		return "", false
	}
	delete(s.items, state)
	if now.After(e.expiresAt) {
		return "", false
	}
	return e.value, true
}

func (s *StateStore) purgeExpiredLocked(now time.Time) {
	for k, e := range s.items {
		if now.After(e.expiresAt) {
			delete(s.items, k)
		}
	}
}

func randomState() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
