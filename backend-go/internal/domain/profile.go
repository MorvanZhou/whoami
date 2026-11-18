package domain

import "context"

// HistoryItem is one historical profile snapshot metadata.
type HistoryItem struct {
	Timestamp    string `json:"timestamp"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
}

// ProfileStore persists user markdown profiles (replaces Tencent COS in the Python backend).
type ProfileStore interface {
	Get(ctx context.Context, userID string) (content string, ok bool, err error)
	Save(ctx context.Context, userID, content string) error
	ListHistory(ctx context.Context, userID string) ([]HistoryItem, error)
	Rollback(ctx context.Context, userID, timestamp string) (bool, error)
}
