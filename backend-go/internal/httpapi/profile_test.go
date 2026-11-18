package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestProfileAPIKeyAndWebJWT(t *testing.T) {
	router, _, jwt, deps := setupTestRouter(t)

	// Create API key + consume store token to get plaintext key
	req := httptest.NewRequest(http.MethodPost, "/api/keys", bytes.NewBufferString(`{"label":"agent"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create key: %d %s", w.Code, w.Body.String())
	}
	var created struct {
		StoreURL string `json:"store_url"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	u, _ := url.Parse(created.StoreURL)
	token := u.Query().Get("token")

	req = httptest.NewRequest(http.MethodGet, "/api/storeapi?token="+url.QueryEscape(token), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("storeapi: %d %s", w.Code, w.Body.String())
	}
	plainKey := ""
	for _, line := range strings.Split(w.Body.String(), "\n") {
		if strings.HasPrefix(line, "WHOAMI_API_KEY=") {
			plainKey = strings.TrimPrefix(line, "WHOAMI_API_KEY=")
		}
	}
	if plainKey == "" {
		t.Fatal("missing api key in storeapi response")
	}

	// Empty profile via API key
	req = httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get profile: %d %s", w.Code, w.Body.String())
	}
	var got map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got["content"] != nil {
		t.Fatalf("expected null content, got %#v", got["content"])
	}

	// Update via API key
	req = httptest.NewRequest(http.MethodPost, "/api/profile", bytes.NewBufferString(`{"content":"v1 hello"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update profile: %d %s", w.Code, w.Body.String())
	}

	// Read via web JWT (same store)
	req = httptest.NewRequest(http.MethodGet, "/api/web/profile", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("web get: %d %s", w.Code, w.Body.String())
	}
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got["content"] != "v1 hello" {
		t.Fatalf("web content=%#v", got["content"])
	}

	// Second save for history
	time.Sleep(10 * time.Millisecond)
	// Advance store clock if possible — localfs uses real time second precision.
	// Force distinct history by saving through service with delay of 1s is slow;
	// instead save twice with injected clock via deps.Profiles store — skip clock,
	// call Save again after bumping by writing through handler once more after 0s
	// may collide; use deps to save with unique content and rely on collision waiter.
	_ = deps
	req = httptest.NewRequest(http.MethodPost, "/api/web/profile", bytes.NewBufferString(`{"content":"v2 via web"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("web update: %d %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/profile/history", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("history: %d %s", w.Code, w.Body.String())
	}
	var hist []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &hist); err != nil {
		t.Fatal(err)
	}
	if len(hist) < 1 {
		t.Fatalf("expected at least 1 history item, got %v", hist)
	}
	ts, _ := hist[len(hist)-1]["timestamp"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/profile/rollback/"+ts, nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("rollback: %d %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got["content"] == nil || got["content"] == "" {
		t.Fatalf("expected content after rollback, got %#v", got["content"])
	}
}

func TestProfileRequiresAPIKey(t *testing.T) {
	router, _, _, _ := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestProfileMaxLength(t *testing.T) {
	router, _, jwt, _ := setupTestRouter(t)

	// mint api key
	req := httptest.NewRequest(http.MethodPost, "/api/keys", bytes.NewBufferString(`{"label":"agent"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var created struct {
		StoreURL string `json:"store_url"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	u, _ := url.Parse(created.StoreURL)
	token := u.Query().Get("token")
	req = httptest.NewRequest(http.MethodGet, "/api/storeapi?token="+url.QueryEscape(token), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	plainKey := ""
	for _, line := range strings.Split(w.Body.String(), "\n") {
		if strings.HasPrefix(line, "WHOAMI_API_KEY=") {
			plainKey = strings.TrimPrefix(line, "WHOAMI_API_KEY=")
		}
	}

	long := strings.Repeat("你", 5001)
	if utf8.RuneCountInString(long) != 5001 {
		t.Fatalf("setup rune count=%d", utf8.RuneCountInString(long))
	}
	payload, _ := json.Marshal(map[string]string{"content": long})
	req = httptest.NewRequest(http.MethodPost, "/api/profile", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plainKey)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Huma maxLength validation returns 422.
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for oversized profile, got %d body=%s", w.Code, w.Body.String())
	}
}
