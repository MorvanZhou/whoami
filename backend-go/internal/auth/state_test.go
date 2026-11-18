package auth

import (
	"testing"
	"time"
)

func TestStateStoreTakeExpires(t *testing.T) {
	s := NewStateStore(20*time.Millisecond, 100)
	key, err := s.Put("/dashboard")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := s.Take(key); !ok || v != "/dashboard" {
		t.Fatalf("take got %q ok=%v", v, ok)
	}
	if _, ok := s.Take(key); ok {
		t.Fatal("expected miss after take")
	}

	key2, err := s.Put("/keys")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Millisecond)
	if _, ok := s.Take(key2); ok {
		t.Fatal("expected expired")
	}
}
