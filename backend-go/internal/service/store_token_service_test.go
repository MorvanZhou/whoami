package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/service"
)

func TestStoreTokenService_ConsumeTable(t *testing.T) {
	ctx := context.Background()
	apiKeyID := uuid.NewString()

	tests := []struct {
		name      string
		setup     func(*mockStoreTokenRepo) (token string)
		wantOK    bool
		wantPlain string
	}{
		{
			name: "consume once",
			setup: func(repo *mockStoreTokenRepo) string {
				svc := service.NewStoreTokenService(repo)
				tok, err := svc.Create(ctx, "u1", "plain-key", apiKeyID)
				if err != nil {
					t.Fatal(err)
				}
				return tok
			},
			wantOK:    true,
			wantPlain: "plain-key",
		},
		{
			name: "missing token",
			setup: func(repo *mockStoreTokenRepo) string {
				return "nope"
			},
			wantOK: false,
		},
		{
			name: "expired token",
			setup: func(repo *mockStoreTokenRepo) string {
				id := apiKeyID
				rec := &domain.StoreApiToken{
					ID:          uuid.NewString(),
					UserID:      "u1",
					Token:       "expired-tok",
					APIKeyPlain: "x",
					APIKeyID:    &id,
					CreatedAt:   time.Now().UTC().Add(-30 * time.Minute),
				}
				_ = repo.Create(ctx, rec)
				return "expired-tok"
			},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockStoreTokenRepo()
			token := tt.setup(repo)
			svc := service.NewStoreTokenService(repo)
			plain, ok, err := svc.Consume(ctx, token)
			if err != nil {
				t.Fatal(err)
			}
			if ok != tt.wantOK {
				t.Fatalf("ok=%v want=%v", ok, tt.wantOK)
			}
			if tt.wantOK && plain != tt.wantPlain {
				t.Fatalf("plain=%q want=%q", plain, tt.wantPlain)
			}
			if tt.wantOK {
				// second consume fails
				_, ok2, err := svc.Consume(ctx, token)
				if err != nil || ok2 {
					t.Fatalf("second consume ok=%v err=%v", ok2, err)
				}
			}
		})
	}
}

func TestStoreTokenService_IsConsumedAndCleanup(t *testing.T) {
	ctx := context.Background()
	repo := newMockStoreTokenRepo()
	svc := service.NewStoreTokenService(repo)

	apiKeyID := uuid.NewString()
	tok, err := svc.Create(ctx, "u1", "plain", apiKeyID)
	if err != nil {
		t.Fatal(err)
	}
	_ = tok

	consumed, err := svc.IsConsumed(ctx, apiKeyID, "u1")
	if err != nil || consumed {
		t.Fatalf("before consume: consumed=%v err=%v", consumed, err)
	}

	_, ok, err := svc.Consume(ctx, tok)
	if err != nil || !ok {
		t.Fatalf("consume: ok=%v err=%v", ok, err)
	}
	consumed, err = svc.IsConsumed(ctx, apiKeyID, "u1")
	if err != nil || !consumed {
		t.Fatalf("after consume: consumed=%v err=%v", consumed, err)
	}

	// Insert expired leftover and cleanup
	id := uuid.NewString()
	_ = repo.Create(ctx, &domain.StoreApiToken{
		ID: uuid.NewString(), UserID: "u1", Token: "old", APIKeyPlain: "p", APIKeyID: &id,
		CreatedAt: time.Now().UTC().Add(-time.Hour),
	})
	n, err := svc.CleanupExpired(ctx)
	if err != nil || n != 1 {
		t.Fatalf("cleanup n=%d err=%v", n, err)
	}
}
