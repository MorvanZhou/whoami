package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/config"
	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/httpapi"
	"whoami/backend-go/internal/service"
	"whoami/backend-go/internal/store/gormdb"
	"whoami/backend-go/internal/store/localfs"
)

func setupTestRouter(t *testing.T) (http.Handler, *domain.User, string, httpapi.Deps) {
	t.Helper()
	db, err := gormdb.Open(":memory:")
	if err != nil {
		t.Fatalf("db: %v", err)
	}
	users := gormdb.NewUserRepo(db)
	keys := gormdb.NewApiKeyRepo(db)
	tokens := gormdb.NewStoreTokenRepo(db)

	cfg := &config.Config{
		Env:           "dev",
		FrontendURL:   "http://localhost:3000",
		APIKeyPrefix:  "wai_",
		JWTSecret:     "test-secret",
		JWTExpireDays: 7,
	}
	tokenMgr := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTExpireDays)
	keySvc := service.NewKeyService(keys, users, cfg.APIKeyPrefix)
	storeSvc := service.NewStoreTokenService(tokens)
	authSvc := service.NewAuthService(users)
	github := auth.NewGitHubOAuth("cid", "csecret", "http://localhost/auth/github/callback")

	profileStore, err := localfs.NewProfileStore(t.TempDir())
	if err != nil {
		t.Fatalf("profile store: %v", err)
	}
	profileSvc := service.NewProfileService(profileStore)

	email := "u@example.com"
	name := "User"
	user, err := authSvc.FindOrCreateUser(context.Background(), "github", "1", &email, &name, nil)
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	jwt, err := tokenMgr.Create(user.ID)
	if err != nil {
		t.Fatalf("jwt: %v", err)
	}

	deps := httpapi.Deps{
		Config:           cfg,
		Tokens:           tokenMgr,
		States:           auth.NewStateStore(10*time.Minute, 1000),
		GitHub:           github,
		Auth:             authSvc,
		Users:            users,
		Keys:             keySvc,
		StoreTokens:      storeSvc,
		Profiles:         profileSvc,
		DisableRateLimit: true, // avoid flaky 429s in suite
	}
	return httpapi.NewRouter(deps), user, jwt, deps
}

func TestKeysAndStoreAPIFlow(t *testing.T) {
	router, _, jwt, _ := setupTestRouter(t)

	body := bytes.NewBufferString(`{"label":"agent"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/keys", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create key: status=%d body=%s", w.Code, w.Body.String())
	}

	var created struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		StoreURL  string `json:"store_url"`
		ExpiresIn int    `json:"expires_in"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.ID == "" || created.StoreURL == "" || created.ExpiresIn != 15*60 {
		t.Fatalf("unexpected create response: %+v", created)
	}

	u, err := url.Parse(created.StoreURL)
	if err != nil {
		t.Fatal(err)
	}
	token := u.Query().Get("token")
	if token == "" {
		t.Fatalf("missing token in store_url: %s", created.StoreURL)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/keys/"+created.ID+"/status", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var st map[string]bool
	_ = json.Unmarshal(w.Body.Bytes(), &st)
	if st["configured"] {
		t.Fatal("expected configured=false before consume")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/storeapi?token="+url.QueryEscape(token), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("storeapi: %d %s", w.Code, w.Body.String())
	}
	text, _ := io.ReadAll(w.Body)
	if !strings.Contains(string(text), "WHOAMI_API_KEY=wai_") {
		t.Fatalf("unexpected body: %s", text)
	}
	if !strings.Contains(string(text), "WHOAMI_ENDPOINT=http://localhost:3000") {
		t.Fatalf("unexpected endpoint: %s", text)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/storeapi?token="+url.QueryEscape(token), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on reuse, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/keys/"+created.ID+"/status", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &st)
	if !st["configured"] {
		t.Fatal("expected configured=true after consume")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/keys", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list: %d %s", w.Code, w.Body.String())
	}
	var list []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil || len(list) != 1 {
		t.Fatalf("list decode: %v len=%d body=%s", err, len(list), w.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/keys/"+created.ID, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete: %d %s", w.Code, w.Body.String())
	}
}

func TestKeysRequiresAuth(t *testing.T) {
	router, _, _, _ := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/keys", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateKeyValidation(t *testing.T) {
	router, _, jwt, _ := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/api/keys", bytes.NewBufferString(`{"label":""}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Huma validation uses 422 Unprocessable Entity (vs Gin binding 400).
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestStoreTokenCleanupExpired(t *testing.T) {
	db, err := gormdb.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	tokens := gormdb.NewStoreTokenRepo(db)
	svc := service.NewStoreTokenService(tokens)
	ctx := context.Background()

	id := uuid.NewString()
	old := &domain.StoreApiToken{
		ID:          uuid.NewString(),
		UserID:      "u1",
		Token:       "old-token",
		APIKeyPlain: "plain",
		APIKeyID:    &id,
		CreatedAt:   time.Now().UTC().Add(-30 * time.Minute),
	}
	if err := tokens.Create(ctx, old); err != nil {
		t.Fatal(err)
	}

	n, err := svc.CleanupExpired(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected 1 cleaned, got %d", n)
	}
}
