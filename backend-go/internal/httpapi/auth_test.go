package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGitHubLoginRedirect(t *testing.T) {
	router, _, _, deps := setupTestRouter(t)
	deps.Config.GitHubClientID = "cid"
	deps.Config.GitHubClientSecret = "csecret"

	req := httptest.NewRequest(http.MethodGet, "/api/auth/github/login?redirect=/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d body=%s", w.Code, w.Body.String())
	}
	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "github.com/login/oauth/authorize") {
		t.Fatalf("unexpected location: %s", loc)
	}
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatal(err)
	}
	if u.Query().Get("state") == "" {
		t.Fatal("missing state")
	}
	if u.Query().Get("client_id") != "cid" {
		t.Fatalf("client_id=%s", u.Query().Get("client_id"))
	}
}

func TestGitHubLoginNotConfigured(t *testing.T) {
	router, _, _, deps := setupTestRouter(t)
	deps.Config.GitHubClientID = ""
	deps.Config.GitHubClientSecret = ""

	req := httptest.NewRequest(http.MethodGet, "/api/auth/github/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestGitHubCallbackAndMeLogout(t *testing.T) {
	router, _, _, deps := setupTestRouter(t)
	deps.Config.GitHubClientID = "cid"
	deps.Config.GitHubClientSecret = "csecret"
	deps.Config.CookieSecure = false

	gh := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/login/oauth/access_token" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"ghtoken","token_type":"bearer","scope":"user:email"}`))
		case r.URL.Path == "/user":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":42,"login":"octocat","name":"The Octocat","email":"octo@example.com","avatar_url":"https://example.com/a.png"}`))
		case r.URL.Path == "/user/emails":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"email":"octo@example.com","primary":true}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(gh.Close)

	deps.GitHub.SetHTTPClient(gh.Client())
	deps.GitHub.SetEndpoints(
		gh.URL+"/login/oauth/authorize",
		gh.URL+"/login/oauth/access_token",
		gh.URL+"/user",
		gh.URL+"/user/emails",
	)

	state, err := deps.States.Put("/dashboard")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/github/callback?code=abc&state="+url.QueryEscape(state), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("callback status=%d body=%s", w.Code, w.Body.String())
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, deps.Config.FrontendURL+"/auth/callback") {
		t.Fatalf("redirect loc=%s", loc)
	}
	if !strings.Contains(loc, "redirect="+url.QueryEscape("/dashboard")) && !strings.Contains(loc, "redirect=%2Fdashboard") {
		t.Fatalf("missing redirect param: %s", loc)
	}
	cookie := w.Result().Header.Get("Set-Cookie")
	if !strings.Contains(cookie, "access_token=") || !strings.Contains(cookie, "HttpOnly") {
		t.Fatalf("missing session cookie: %s", cookie)
	}

	// Extract JWT from Set-Cookie for /auth/me
	jwt := extractCookie(cookie, "access_token")
	if jwt == "" {
		t.Fatalf("empty jwt cookie: %s", cookie)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("me status=%d body=%s", w.Code, w.Body.String())
	}
	var me map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &me); err != nil {
		t.Fatal(err)
	}
	if me["provider"] != "github" {
		t.Fatalf("me=%v", me)
	}
	if me["email"] != "octo@example.com" {
		t.Fatalf("email=%v", me["email"])
	}

	req = httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout=%d", w.Code)
	}
	if !strings.Contains(w.Header().Get("Set-Cookie"), "access_token=") {
		t.Fatalf("expected cookie clear: %s", w.Header().Get("Set-Cookie"))
	}
}

func TestGitHubCallbackInvalidState(t *testing.T) {
	router, _, _, _ := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/auth/github/callback?code=abc&state=nope", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func extractCookie(setCookie, name string) string {
	parts := strings.Split(setCookie, ";")
	if len(parts) == 0 {
		return ""
	}
	kv := strings.TrimSpace(parts[0])
	prefix := name + "="
	if strings.HasPrefix(kv, prefix) {
		return strings.TrimPrefix(kv, prefix)
	}
	return ""
}
