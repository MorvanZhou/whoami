package localfs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestStore(t *testing.T) (*ProfileStore, *time.Time) {
	t.Helper()
	root := t.TempDir()
	store, err := NewProfileStore(root)
	if err != nil {
		t.Fatalf("NewProfileStore: %v", err)
	}
	base := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	now := base
	store.now = func() time.Time { return now }
	return store, &now
}

func TestProfileStore_SaveGetHistoryRollback(t *testing.T) {
	ctx := context.Background()
	store, now := newTestStore(t)
	userID := "user-123"

	content, ok, err := store.Get(ctx, userID)
	if err != nil || ok || content != "" {
		t.Fatalf("empty get: content=%q ok=%v err=%v", content, ok, err)
	}

	if err := store.Save(ctx, userID, "v1"); err != nil {
		t.Fatalf("save v1: %v", err)
	}
	content, ok, err = store.Get(ctx, userID)
	if err != nil || !ok || content != "v1" {
		t.Fatalf("get v1: content=%q ok=%v err=%v", content, ok, err)
	}

	*now = now.Add(time.Second)
	if err := store.Save(ctx, userID, "v2"); err != nil {
		t.Fatalf("save v2: %v", err)
	}
	*now = now.Add(time.Second)
	if err := store.Save(ctx, userID, "v3"); err != nil {
		t.Fatalf("save v3: %v", err)
	}
	*now = now.Add(time.Second)
	if err := store.Save(ctx, userID, "v4"); err != nil {
		t.Fatalf("save v4: %v", err)
	}

	content, ok, err = store.Get(ctx, userID)
	if err != nil || !ok || content != "v4" {
		t.Fatalf("get v4: content=%q ok=%v err=%v", content, ok, err)
	}

	hist, err := store.ListHistory(ctx, userID)
	if err != nil {
		t.Fatalf("list history: %v", err)
	}
	if len(hist) != 3 {
		t.Fatalf("want 3 history items after cleanup, got %d: %+v", len(hist), hist)
	}

	// After saves v2/v3/v4, history should be backups of v1,v2,v3 (oldest v1 dropped? 
	// Save v1: no backup. Save v2: backup v1. Save v3: backup v2. Save v4: backup v3.
	// cleanup keeps latest 3 => v1,v2,v3 backups. Good.
	if hist[0].Timestamp < hist[1].Timestamp {
		t.Fatalf("history should be newest-first: %+v", hist)
	}

	oldestKept := hist[len(hist)-1].Timestamp
	data, err := os.ReadFile(filepath.Join(store.rootDir, userID, "history", oldestKept+".md"))
	if err != nil {
		t.Fatalf("read oldest history: %v", err)
	}
	want := string(data)

	*now = now.Add(time.Second)
	okRB, err := store.Rollback(ctx, userID, oldestKept)
	if err != nil || !okRB {
		t.Fatalf("rollback: ok=%v err=%v", okRB, err)
	}

	content, ok, err = store.Get(ctx, userID)
	if err != nil || !ok || content != want {
		t.Fatalf("get after rollback: content=%q want=%q ok=%v err=%v", content, want, ok, err)
	}

	hist2, err := store.ListHistory(ctx, userID)
	if err != nil {
		t.Fatalf("list after rollback: %v", err)
	}
	if len(hist2) > 3 {
		t.Fatalf("history should stay <= 3, got %d", len(hist2))
	}
}

func TestProfileStore_RejectsUnsafeIDs(t *testing.T) {
	ctx := context.Background()
	store, err := NewProfileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := store.Get(ctx, "../etc"); err == nil {
		t.Fatal("expected invalid userID error")
	}
	if err := store.Save(ctx, "a/b", "x"); err == nil {
		t.Fatal("expected invalid userID error")
	}
	if _, err := store.Rollback(ctx, "user1", "../x"); err == nil {
		t.Fatal("expected invalid timestamp error")
	}
}

func TestProfileStore_RollbackMissing(t *testing.T) {
	ctx := context.Background()
	store, err := NewProfileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ok, err := store.Rollback(ctx, "user1", "20200101T000000Z")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ok {
		t.Fatal("expected ok=false for missing history")
	}
}
