package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/service"
	"whoami/backend-go/internal/store/gormdb"
)

func strPtr(s string) *string { return &s }

func mustCreateAPIKey(t *testing.T, ctx context.Context, keys domain.ApiKeyRepository, userID, label string) *domain.ApiKey {
	t.Helper()
	lbl := label
	key := &domain.ApiKey{
		ID:        uuid.NewString(),
		UserID:    userID,
		KeyHash:   "hash-" + uuid.NewString(),
		KeyPrefix: "wai_test",
		KeySuffix: "abcd",
		Label:     &lbl,
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}
	if err := keys.Create(ctx, key); err != nil {
		t.Fatalf("create api key: %v", err)
	}
	return key
}

func mustCreateStoreToken(t *testing.T, ctx context.Context, tokens domain.StoreTokenRepository, userID, apiKeyID, plain string) *domain.StoreApiToken {
	t.Helper()
	id := apiKeyID
	tok := &domain.StoreApiToken{
		ID:          uuid.NewString(),
		UserID:      userID,
		Token:       "tok-" + uuid.NewString(),
		APIKeyPlain: plain,
		APIKeyID:    &id,
		IsUsed:      false,
		CreatedAt:   time.Now().UTC(),
	}
	if err := tokens.Create(ctx, tok); err != nil {
		t.Fatalf("create store token: %v", err)
	}
	return tok
}

func TestFindOrCreateUser(t *testing.T) {
	ctx := context.Background()
	db, err := gormdb.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	users := gormdb.NewUserRepo(db)
	svc := service.NewAuthService(users)

	email := strPtr("a@example.com")
	name := strPtr("Alice")

	u1, err := svc.FindOrCreateUser(ctx, "github", "gh-1", email, name, nil)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if u1.ID == "" || u1.Provider != "github" {
		t.Fatalf("unexpected user: %+v", u1)
	}

	// Same provider identity → same user, profile updated.
	newName := strPtr("Alice Updated")
	u2, err := svc.FindOrCreateUser(ctx, "github", "gh-1", email, newName, nil)
	if err != nil {
		t.Fatalf("login again: %v", err)
	}
	if u2.ID != u1.ID {
		t.Fatalf("expected same id, got %s vs %s", u1.ID, u2.ID)
	}
	if u2.Name == nil || *u2.Name != "Alice Updated" {
		t.Fatalf("expected updated name, got %+v", u2.Name)
	}

	// Different provider, same email → link to existing (Python behavior).
	u3, err := svc.FindOrCreateUser(ctx, "google", "gg-1", email, strPtr("G"), nil)
	if err != nil {
		t.Fatalf("link by email: %v", err)
	}
	if u3.ID != u1.ID {
		t.Fatalf("expected email-linked same user, got %s vs %s", u3.ID, u1.ID)
	}
	// Original provider fields are intentionally NOT overwritten in Python.
	if u3.Provider != "github" {
		t.Fatalf("expected original provider github, got %s", u3.Provider)
	}

	// Brand-new user with no overlapping email.
	u4, err := svc.FindOrCreateUser(ctx, "github", "gh-2", strPtr("b@example.com"), strPtr("Bob"), nil)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}
	if u4.ID == u1.ID {
		t.Fatal("expected a different user id")
	}
}

func TestUserRepoCRUD(t *testing.T) {
	ctx := context.Background()
	db, err := gormdb.Open("file:phase2?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	users := gormdb.NewUserRepo(db)
	keys := gormdb.NewApiKeyRepo(db)
	tokens := gormdb.NewStoreTokenRepo(db)

	svc := service.NewAuthService(users)
	u, err := svc.FindOrCreateUser(ctx, "github", "x", strPtr("x@ex.com"), strPtr("X"), nil)
	if err != nil {
		t.Fatal(err)
	}

	label := "default"
	key := mustCreateAPIKey(t, ctx, keys, u.ID, label)
	got, err := keys.FindActiveByHash(ctx, key.KeyHash)
	if err != nil || got == nil || got.ID != key.ID {
		t.Fatalf("find by hash: got=%+v err=%v", got, err)
	}

	list, err := keys.ListActiveByUser(ctx, u.ID)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v len=%d", err, len(list))
	}

	tok := mustCreateStoreToken(t, ctx, tokens, u.ID, key.ID, "plain-key")
	found, err := tokens.FindUnusedByToken(ctx, tok.Token)
	if err != nil || found == nil {
		t.Fatalf("find token: %v", err)
	}
	exists, err := tokens.ExistsByAPIKeyID(ctx, key.ID, u.ID)
	if err != nil || !exists {
		t.Fatalf("exists: %v %v", exists, err)
	}
	if err := tokens.Delete(ctx, tok.ID); err != nil {
		t.Fatal(err)
	}
	exists, err = tokens.ExistsByAPIKeyID(ctx, key.ID, u.ID)
	if err != nil || exists {
		t.Fatalf("expected consumed/deleted token, exists=%v err=%v", exists, err)
	}

	ok, err := keys.DeleteByIDForUser(ctx, key.ID, u.ID)
	if err != nil || !ok {
		t.Fatalf("delete key: ok=%v err=%v", ok, err)
	}
}
