package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"whoami/backend-go/internal/service"
)

func TestProfileService_GetTable(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(*mockProfileStore)
		userID    string
		wantNil   bool
		wantVal   string
		wantErr   bool
	}{
		{
			name:    "missing profile returns nil content",
			userID:  "u1",
			wantNil: true,
		},
		{
			name: "existing profile",
			setup: func(m *mockProfileStore) {
				_ = m.Save(ctx, "u1", "hello")
			},
			userID:  "u1",
			wantVal: "hello",
		},
		{
			name: "store error",
			setup: func(m *mockProfileStore) {
				m.getErr = errors.New("boom")
			},
			userID:  "u1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockProfileStore()
			if tt.setup != nil {
				tt.setup(store)
			}
			svc := service.NewProfileService(store)
			got, err := svc.Get(ctx, tt.userID)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Fatalf("want nil, got %v", *got)
				}
				return
			}
			if got == nil || *got != tt.wantVal {
				t.Fatalf("got=%v want=%q", got, tt.wantVal)
			}
		})
	}
}

func TestProfileService_SaveTable(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		content string
		saveErr error
		wantErr string
	}{
		{name: "ok", content: "v1"},
		{name: "too long", content: strings.Repeat("あ", service.ProfileMaxLength+1), wantErr: "exceeds"},
		{name: "store failure", content: "x", saveErr: errors.New("disk full"), wantErr: "disk full"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMockProfileStore()
			store.saveErr = tt.saveErr
			svc := service.NewProfileService(store)
			err := svc.Save(ctx, "u1", tt.content)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				got, ok, _ := store.Get(ctx, "u1")
				if !ok || got != tt.content {
					t.Fatalf("stored=%q ok=%v", got, ok)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err=%v want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestProfileService_RollbackWithMock(t *testing.T) {
	ctx := context.Background()
	store := newMockProfileStore()
	svc := service.NewProfileService(store)

	if err := svc.Save(ctx, "u1", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Save(ctx, "u1", "v2"); err != nil {
		t.Fatal(err)
	}
	hist, err := svc.ListHistory(ctx, "u1")
	if err != nil || len(hist) == 0 {
		t.Fatalf("history=%v err=%v", hist, err)
	}
	ok, err := svc.Rollback(ctx, "u1", hist[0].Timestamp)
	if err != nil || !ok {
		t.Fatalf("rollback ok=%v err=%v", ok, err)
	}
	content, err := svc.Get(ctx, "u1")
	if err != nil || content == nil || *content != "v1" {
		t.Fatalf("after rollback content=%v err=%v", content, err)
	}
}
