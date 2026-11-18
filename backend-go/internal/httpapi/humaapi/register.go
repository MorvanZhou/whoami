package humaapi

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/config"
	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/httpapi/middleware"
	"whoami/backend-go/internal/service"
)

// Deps for Huma JSON API operations.
type Deps struct {
	Config      *config.Config
	Tokens      *auth.TokenManager
	Users       domain.UserRepository
	Keys        *service.KeyService
	StoreTokens *service.StoreTokenService
	Profiles    *service.ProfileService
}

// Register mounts Huma operations on the Gin /api group and returns the API (for docs/tests).
// OAuth redirects and plaintext /storeapi stay on pure Gin.
func Register(r *gin.Engine, apiGroup *gin.RouterGroup, d Deps) huma.API {
	cfg := huma.DefaultConfig("whoami API", "0.1.0")
	cfg.Info.Description = "Cross-AI user identity profile sync service (Go practice backend). Hybrid: JSON routes via Huma; OAuth/storeapi remain Gin."
	// NewWithGroup mounts under /api; Servers tells the docs UI to fetch /api/openapi.yaml
	// (otherwise Elements requests /openapi.yaml at the site root → 404).
	cfg.Servers = []*huma.Server{{URL: "/api"}}
	if d.Config.Env != "dev" {
		cfg.DocsPath = ""
		cfg.SchemasPath = ""
	}
	if cfg.Components == nil {
		cfg.Components = &huma.Components{}
	}
	if cfg.Components.SecuritySchemes == nil {
		cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	cfg.Components.SecuritySchemes["bearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT from OAuth cookie or Authorization: Bearer <jwt>",
	}
	cfg.Components.SecuritySchemes["apiKeyAuth"] = &huma.SecurityScheme{
		Type:        "http",
		Scheme:      "bearer",
		Description: "whoami API key as Authorization: Bearer <wai_...>",
	}

	api := humagin.NewWithGroup(r, apiGroup, cfg)

	jwtMW := middleware.HumaRequireJWT(api, d.Tokens, d.Users)
	apiKeyMW := middleware.HumaRequireAPIKey(api, d.Keys)

	registerAuthMe(api, jwtMW)
	registerKeys(api, d, jwtMW)
	registerProfiles(api, d, jwtMW, apiKeyMW)

	return api
}
