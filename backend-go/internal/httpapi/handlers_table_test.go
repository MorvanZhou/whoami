package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateKey_TableDriven(t *testing.T) {
	router, _, jwt, _ := setupTestRouter(t)

	tests := []struct {
		name       string
		authHeader string
		body       string
		wantStatus int
	}{
		{
			name:       "ok",
			authHeader: "Bearer " + jwt,
			body:       `{"label":"agent"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing auth",
			body:       `{"label":"agent"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "empty label",
			authHeader: "Bearer " + jwt,
			body:       `{"label":""}`,
			wantStatus: http.StatusUnprocessableEntity, // Huma validation
		},
		{
			name:       "label too long",
			authHeader: "Bearer " + jwt,
			body:       `{"label":"abcdefghijklmnop"}`, // 16 chars > 15
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid json",
			authHeader: "Bearer " + jwt,
			body:       `{`,
			wantStatus: http.StatusBadRequest, // malformed JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/keys", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status=%d body=%s want=%d", w.Code, w.Body.String(), tt.wantStatus)
			}
			if tt.wantStatus == http.StatusCreated {
				var resp map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatal(err)
				}
				if resp["store_url"] == nil || resp["id"] == nil {
					t.Fatalf("missing fields: %v", resp)
				}
			}
		})
	}
}

func TestAuthMe_TableDriven(t *testing.T) {
	router, user, jwt, _ := setupTestRouter(t)

	tests := []struct {
		name       string
		authHeader string
		cookie     string
		wantStatus int
		wantID     string
	}{
		{
			name:       "bearer jwt",
			authHeader: "Bearer " + jwt,
			wantStatus: http.StatusOK,
			wantID:     user.ID,
		},
		{
			name:       "cookie jwt",
			cookie:     "access_token=" + jwt,
			wantStatus: http.StatusOK,
			wantID:     user.ID,
		},
		{
			name:       "missing",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "garbage bearer",
			authHeader: "Bearer not-a-jwt",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.cookie != "" {
				req.Header.Set("Cookie", tt.cookie)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tt.wantStatus == http.StatusOK {
				var me map[string]any
				_ = json.Unmarshal(w.Body.Bytes(), &me)
				if me["id"] != tt.wantID {
					t.Fatalf("id=%v want=%s", me["id"], tt.wantID)
				}
			}
		})
	}
}

func TestStoreAPI_TableDriven(t *testing.T) {
	router, _, jwt, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/keys", bytes.NewBufferString(`{"label":"t"}`))
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
	u, err := url.Parse(created.StoreURL)
	if err != nil {
		t.Fatal(err)
	}
	token := u.Query().Get("token")
	if token == "" {
		t.Fatalf("no token in %s", created.StoreURL)
	}

	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantSubstr string
	}{
		{name: "missing token", query: "", wantStatus: http.StatusBadRequest},
		{name: "bad token", query: "?token=nope", wantStatus: http.StatusNotFound},
		{name: "good token", query: "?token=" + url.QueryEscape(token), wantStatus: http.StatusOK, wantSubstr: "WHOAMI_API_KEY="},
		{name: "reuse token", query: "?token=" + url.QueryEscape(token), wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/storeapi"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status=%d body=%s want=%d", w.Code, w.Body.String(), tt.wantStatus)
			}
			if tt.wantSubstr != "" && !strings.Contains(w.Body.String(), tt.wantSubstr) {
				t.Fatalf("body=%s missing %s", w.Body.String(), tt.wantSubstr)
			}
		})
	}
}

func TestProfileAuth_TableDriven(t *testing.T) {
	router, _, jwt, _ := setupTestRouter(t)

	tests := []struct {
		name       string
		method     string
		path       string
		auth       string
		wantStatus int
	}{
		{name: "profile needs api key", method: http.MethodGet, path: "/api/profile", wantStatus: http.StatusUnauthorized},
		{name: "profile rejects jwt", method: http.MethodGet, path: "/api/profile", auth: "Bearer " + jwt, wantStatus: http.StatusUnauthorized},
		{name: "web profile needs jwt", method: http.MethodGet, path: "/api/web/profile", wantStatus: http.StatusUnauthorized},
		{name: "web profile accepts jwt", method: http.MethodGet, path: "/api/web/profile", auth: "Bearer " + jwt, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status=%d body=%s want=%d", w.Code, w.Body.String(), tt.wantStatus)
			}
		})
	}
}
