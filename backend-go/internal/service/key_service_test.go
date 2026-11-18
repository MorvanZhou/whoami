package service_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/service"
)

func TestKeyService_CreateListRevokeAuthenticate(t *testing.T) {
	ctx := context.Background()
	users := newMockUserRepo()
	keys := newMockAPIKeyRepo()
	svc := service.NewKeyService(keys, users, "wai_")

	user := &domain.User{ID: "user-1", Provider: "github", ProviderID: "1", IsActive: true, CreatedAt: time.Now().UTC()}
	_ = users.Create(ctx, user)

	plain, key, err := svc.Create(ctx, user.ID, "agent")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(plain, "wai_") || key == nil || key.UserID != user.ID {
		t.Fatalf("plain=%s key=%+v", plain, key)
	}
	if key.KeyPrefix != plain[:8] || key.KeySuffix != plain[len(plain)-4:] {
		t.Fatalf("prefix/suffix mismatch: %+v plain=%s", key, plain)
	}

	list, err := svc.List(ctx, user.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("list=%v err=%v", list, err)
	}

	gotUser, err := svc.Authenticate(ctx, plain)
	if err != nil || gotUser == nil || gotUser.ID != user.ID {
		t.Fatalf("auth user=%v err=%v", gotUser, err)
	}

	// Wrong key
	bad, err := svc.Authenticate(ctx, "wai_"+strings.Repeat("0", 48))
	if err != nil || bad != nil {
		t.Fatalf("expected nil user for bad key, got=%v err=%v", bad, err)
	}

	ok, err := svc.Revoke(ctx, key.ID, user.ID)
	if err != nil || !ok {
		t.Fatalf("revoke ok=%v err=%v", ok, err)
	}
	list, err = svc.List(ctx, user.ID)
	if err != nil || len(list) != 0 {
		t.Fatalf("after revoke list=%v err=%v", list, err)
	}
}

func TestKeyService_AuthenticateTable(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(users *mockUserRepo, keys *mockAPIKeyRepo) (plain string)
		wantUserID string
		wantNil    bool
	}{
		{
			name: "active user",
			setup: func(users *mockUserRepo, keys *mockAPIKeyRepo) string {
				u := &domain.User{ID: "u1", IsActive: true}
				_ = users.Create(ctx, u)
				plain, key, err := service.NewKeyService(keys, users, "wai_").Create(ctx, "u1", "l")
				if err != nil {
					t.Fatal(err)
				}
				_ = key
				return plain
			},
			wantUserID: "u1",
		},
		{
			name: "inactive user",
			setup: func(users *mockUserRepo, keys *mockAPIKeyRepo) string {
				u := &domain.User{ID: "u2", IsActive: false}
				_ = users.Create(ctx, u)
				plain, _, err := service.NewKeyService(keys, users, "wai_").Create(ctx, "u2", "l")
				if err != nil {
					t.Fatal(err)
				}
				return plain
			},
			wantNil: true,
		},
		{
			name: "unknown key",
			setup: func(users *mockUserRepo, keys *mockAPIKeyRepo) string {
				return "wai_deadbeef"
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := newMockUserRepo()
			keys := newMockAPIKeyRepo()
			plain := tt.setup(users, keys)
			svc := service.NewKeyService(keys, users, "wai_")
			u, err := svc.Authenticate(ctx, plain)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantNil {
				if u != nil {
					t.Fatalf("want nil, got %+v", u)
				}
				return
			}
			if u == nil || u.ID != tt.wantUserID {
				t.Fatalf("got=%+v want=%s", u, tt.wantUserID)
			}
		})
	}
}

func TestHashStable(t *testing.T) {
	// Ensure Authenticate hash matches Create hash (sha256 hex).
	plain := "wai_abc"
	sum := sha256.Sum256([]byte(plain))
	want := hex.EncodeToString(sum[:])

	users := newMockUserRepo()
	keys := newMockAPIKeyRepo()
	_ = users.Create(context.Background(), &domain.User{ID: "u", IsActive: true})
	key := &domain.ApiKey{
		ID: "k1", UserID: "u", KeyHash: want, KeyPrefix: "wai_abc", KeySuffix: "_abc", IsActive: true,
		CreatedAt: time.Now().UTC(),
	}
	_ = keys.Create(context.Background(), key)
	svc := service.NewKeyService(keys, users, "wai_")
	u, err := svc.Authenticate(context.Background(), plain)
	if err != nil || u == nil {
		t.Fatalf("user=%v err=%v", u, err)
	}
}
