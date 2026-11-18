package middleware

import (
	"fmt"
	"testing"
	"time"
)

func TestIPRateLimiter_CleanupOnInsertWhenFull(t *testing.T) {
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	l := NewIPRateLimiter(60, 60)
	l.maxEntries = 2
	l.idleTTL = 10 * time.Minute
	l.now = func() time.Time { return now }

	_ = l.get("1.1.1.1")
	_ = l.get("2.2.2.2")
	if l.Len() != 2 {
		t.Fatalf("want 2 entries, got %d", l.Len())
	}

	now = now.Add(5 * time.Minute)
	_ = l.get("1.1.1.1") // refresh — still under max, no cleanup

	now = now.Add(6 * time.Minute) // 2.2.2.2 idle 11m; 1.1.1.1 idle 6m
	_ = l.get("3.3.3.3")           // triggers cleanup then insert

	if l.Len() != 2 {
		t.Fatalf("want 2 entries after cleanup+insert, got %d", l.Len())
	}
	if _, ok := l.limiters["1.1.1.1"]; !ok {
		t.Fatal("expected 1.1.1.1 to remain")
	}
	if _, ok := l.limiters["2.2.2.2"]; ok {
		t.Fatal("expected 2.2.2.2 to be removed")
	}
	if _, ok := l.limiters["3.3.3.3"]; !ok {
		t.Fatal("expected 3.3.3.3 to be inserted")
	}
}

func TestIPRateLimiter_NoCleanupBelowMaxEntries(t *testing.T) {
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	l := NewIPRateLimiter(60, 60)
	l.maxEntries = 100
	l.idleTTL = time.Minute
	l.now = func() time.Time { return now }

	for i := 0; i < 10; i++ {
		_ = l.get(fmt.Sprintf("10.0.0.%d", i))
	}
	now = now.Add(2 * time.Minute)
	_ = l.get("10.0.0.99")
	if l.Len() != 11 {
		t.Fatalf("idle entries should stay until map is full, got %d", l.Len())
	}
}
