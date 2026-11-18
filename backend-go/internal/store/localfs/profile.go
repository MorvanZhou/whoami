package localfs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"whoami/backend-go/internal/domain"
)

var (
	safeIDRe        = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	safeTimestampRe = regexp.MustCompile(`^[0-9]{8}T[0-9]{6}Z$`)
)

const maxHistory = 3

// ProfileStore stores profiles under rootDir/{userID}/...
// Layout mirrors the Python COS keys:
//
//	{userID}/current.md
//	{userID}/history/{timestamp}.md
type ProfileStore struct {
	rootDir string
	now     func() time.Time

	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

// NewProfileStore creates a local filesystem profile store.
func NewProfileStore(rootDir string) (*ProfileStore, error) {
	if rootDir == "" {
		return nil, errors.New("profile rootDir is empty")
	}
	abs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, err
	}
	return &ProfileStore{
		rootDir: abs,
		now:     func() time.Time { return time.Now().UTC() },
		locks:   make(map[string]*sync.Mutex),
	}, nil
}

func (s *ProfileStore) userLock(userID string) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()
	if l, ok := s.locks[userID]; ok {
		return l
	}
	l := &sync.Mutex{}
	s.locks[userID] = l
	return l
}

func validateUserID(userID string) error {
	if !safeIDRe.MatchString(userID) {
		return fmt.Errorf("invalid userID")
	}
	return nil
}

func validateTimestamp(ts string) error {
	if !safeTimestampRe.MatchString(ts) {
		return fmt.Errorf("invalid timestamp")
	}
	return nil
}

func (s *ProfileStore) userDir(userID string) string {
	return filepath.Join(s.rootDir, userID)
}

func (s *ProfileStore) currentPath(userID string) string {
	return filepath.Join(s.userDir(userID), "current.md")
}

func (s *ProfileStore) historyDir(userID string) string {
	return filepath.Join(s.userDir(userID), "history")
}

func (s *ProfileStore) historyPath(userID, timestamp string) string {
	return filepath.Join(s.historyDir(userID), timestamp+".md")
}

// Get returns the current profile content.
func (s *ProfileStore) Get(_ context.Context, userID string) (string, bool, error) {
	if err := validateUserID(userID); err != nil {
		return "", false, err
	}
	lock := s.userLock(userID)
	lock.Lock()
	defer lock.Unlock()

	data, err := os.ReadFile(s.currentPath(userID))
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return string(data), true, nil
}

// Save writes content to current.md after backing up the previous version.
func (s *ProfileStore) Save(_ context.Context, userID, content string) error {
	if err := validateUserID(userID); err != nil {
		return err
	}
	lock := s.userLock(userID)
	lock.Lock()
	defer lock.Unlock()

	if err := os.MkdirAll(s.historyDir(userID), 0o755); err != nil {
		return err
	}

	if err := s.backupCurrentLocked(userID); err != nil {
		return err
	}
	if err := atomicWriteFile(s.currentPath(userID), []byte(content)); err != nil {
		return err
	}
	return s.cleanupHistoryLocked(userID)
}

// ListHistory returns history items newest-first.
func (s *ProfileStore) ListHistory(_ context.Context, userID string) ([]domain.HistoryItem, error) {
	if err := validateUserID(userID); err != nil {
		return nil, err
	}
	lock := s.userLock(userID)
	lock.Lock()
	defer lock.Unlock()

	entries, err := listHistoryFiles(s.historyDir(userID))
	if err != nil {
		return nil, err
	}

	items := make([]domain.HistoryItem, 0, len(entries))
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		ts := trimMD(e.Name())
		info, err := e.Info()
		if err != nil {
			continue
		}
		items = append(items, domain.HistoryItem{
			Timestamp:    ts,
			Size:         info.Size(),
			LastModified: info.ModTime().UTC().Format(time.RFC3339),
		})
	}
	return items, nil
}

// Rollback restores a history version to current (after backing up current).
func (s *ProfileStore) Rollback(_ context.Context, userID, timestamp string) (bool, error) {
	if err := validateUserID(userID); err != nil {
		return false, err
	}
	if err := validateTimestamp(timestamp); err != nil {
		return false, err
	}

	lock := s.userLock(userID)
	lock.Lock()
	defer lock.Unlock()

	histPath := s.historyPath(userID, timestamp)
	if _, err := os.Stat(histPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if err := os.MkdirAll(s.historyDir(userID), 0o755); err != nil {
		return false, err
	}
	if err := s.backupCurrentLocked(userID); err != nil {
		return false, err
	}

	data, err := os.ReadFile(histPath)
	if err != nil {
		return false, err
	}
	if err := atomicWriteFile(s.currentPath(userID), data); err != nil {
		return false, err
	}
	if err := s.cleanupHistoryLocked(userID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *ProfileStore) backupCurrentLocked(userID string) error {
	cur := s.currentPath(userID)
	if _, err := os.Stat(cur); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Timestamp format matches Python COS keys and rollback validation.
	// If two backups collide in the same second, advance the clock by 1s.
	ts := s.now().Format("20060102T150405Z")
	dst := s.historyPath(userID, ts)
	for i := 0; i < 1000; i++ {
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			return copyFile(cur, dst)
		}
		ts = s.now().Add(time.Duration(i+1) * time.Second).Format("20060102T150405Z")
		dst = s.historyPath(userID, ts)
	}
	return fmt.Errorf("could not allocate unique history timestamp for user %s", userID)
}

func (s *ProfileStore) cleanupHistoryLocked(userID string) error {
	entries, err := listHistoryFiles(s.historyDir(userID))
	if err != nil {
		return err
	}
	if len(entries) <= maxHistory {
		return nil
	}
	toDelete := entries[:len(entries)-maxHistory]
	for _, e := range toDelete {
		if err := os.Remove(filepath.Join(s.historyDir(userID), e.Name())); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func listHistoryFiles(dir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	filtered := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".md" {
			continue
		}
		ts := trimMD(name)
		if !safeTimestampRe.MatchString(ts) {
			continue
		}
		filtered = append(filtered, e)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name() < filtered[j].Name()
	})
	return filtered, nil
}

func trimMD(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}

func atomicWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".current-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp, err := os.CreateTemp(filepath.Dir(dst), ".hist-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := io.Copy(tmp, in); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, dst)
}
