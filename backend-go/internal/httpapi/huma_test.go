package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"whoami/backend-go/internal/httpapi"
)

func mustNewRouter(t *testing.T, deps httpapi.Deps) http.Handler {
	t.Helper()
	return httpapi.NewRouter(deps)
}

func TestHumaOpenAPIAndDocs(t *testing.T) {
	router, _, _, deps := setupTestRouter(t)
	if deps.Config.Env != "dev" {
		t.Fatal("tests expect ENV=dev for docs")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		// Some huma versions serve openapi.yaml / openapi.json at different paths.
		req = httptest.NewRequest(http.MethodGet, "/api/openapi.yaml", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("openapi status=%d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "whoami") && !strings.Contains(body, "auth/me") && !strings.Contains(body, "get-auth-me") {
		// JSON openapi should mention paths or title
		var doc map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
			// yaml
			if !strings.Contains(body, "/auth/me") && !strings.Contains(body, "whoami API") {
				t.Fatalf("unexpected openapi body: %s", truncate(body, 500))
			}
		} else {
			info, _ := doc["info"].(map[string]any)
			if info == nil || info["title"] == nil {
				t.Fatalf("missing info.title in openapi json: %v", doc["info"])
			}
		}
	}

	req = httptest.NewRequest(http.MethodGet, "/api/docs", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("docs status=%d body=%s", w.Code, w.Body.String())
	}
	docsBody := w.Body.String()
	if !strings.Contains(docsBody, `apiDescriptionUrl="/api/openapi.yaml"`) {
		t.Fatalf("docs must point OpenAPI under /api prefix, got: %s", truncate(docsBody, 800))
	}
	if !strings.Contains(docsBody, "html") && !strings.Contains(strings.ToLower(docsBody), "swagger") && !strings.Contains(docsBody, "openapi") {
		if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") && !strings.Contains(ct, "json") {
			t.Fatalf("unexpected docs content-type=%s body=%s", ct, truncate(docsBody, 200))
		}
	}
}

func TestHumaDocsDisabledInProd(t *testing.T) {
	_, _, _, deps := setupTestRouter(t)
	deps.Config.Env = "prod"
	r := mustNewRouter(t, deps)

	req := httptest.NewRequest(http.MethodGet, "/api/docs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK && strings.Contains(strings.ToLower(w.Body.String()), "<html") {
		t.Fatalf("expected docs disabled in prod, got 200 html")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
