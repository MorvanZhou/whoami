package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"whoami/backend-go/internal/httpapi"
	"whoami/backend-go/internal/httpapi/middleware"
)

func TestCORSHeaders(t *testing.T) {
	router, _, _, deps := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	req.Header.Set("Origin", deps.Config.FrontendURL)
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != deps.Config.FrontendURL {
		t.Fatalf("Allow-Origin=%q want %q status=%d", got, deps.Config.FrontendURL, w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("expected credentials true, headers=%v", w.Header())
	}
}

func TestAPIHealth(t *testing.T) {
	router, _, _, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != `{"status":"ok"}` {
		t.Fatalf("/api/health => %d %s", w.Code, w.Body.String())
	}
}

func TestNotFoundDetail(t *testing.T) {
	router, _, _, _ := setupTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d", w.Code)
	}
	if body := w.Body.String(); body != `{"detail":"Not Found"}` {
		t.Fatalf("body=%s", body)
	}
}

func TestRateLimitExceeded(t *testing.T) {
	_, _, _, deps := setupTestRouter(t)
	// Rebuild router with a tiny login limiter enabled.
	deps.DisableRateLimit = false
	deps.GlobalLimiter = middleware.NewIPRateLimiter(60, 60)
	deps.LoginLimiter = middleware.NewIPRateLimiter(2, 2) // 2/min burst 2
	deps.StoreLimiter = middleware.NewIPRateLimiter(10, 10)
	r := httpapi.NewRouter(deps)

	deps.Config.GitHubClientID = "cid"
	deps.Config.GitHubClientSecret = "sec"

	var last int
	for i := 0; i < 4; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/github/login", nil)
		req.RemoteAddr = "203.0.113.10:1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		last = w.Code
		if i < 2 && w.Code != http.StatusFound && w.Code != http.StatusOK {
			// redirect expected
			if w.Code == http.StatusTooManyRequests {
				t.Fatalf("hit rate limit too early at i=%d", i)
			}
		}
	}
	if last != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after bursting login limiter, last=%d", last)
	}
}
