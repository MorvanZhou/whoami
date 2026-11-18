package httpapi

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/config"
	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/httpapi/handlers"
	"whoami/backend-go/internal/httpapi/humaapi"
	"whoami/backend-go/internal/httpapi/middleware"
	"whoami/backend-go/internal/service"
)

// Deps are runtime dependencies for HTTP routes.
type Deps struct {
	Config      *config.Config
	Tokens      *auth.TokenManager
	States      *auth.StateStore
	GitHub      *auth.GitHubOAuth
	Auth        *service.AuthService
	Users       domain.UserRepository
	Keys        *service.KeyService
	StoreTokens *service.StoreTokenService
	Profiles    *service.ProfileService

	// Optional limiters (nil → defaults). Set DisableRateLimit in tests.
	DisableRateLimit bool
	GlobalLimiter    *middleware.IPRateLimiter
	LoginLimiter     *middleware.IPRateLimiter
	StoreLimiter     *middleware.IPRateLimiter
}

// NewRouter builds the Gin engine (hybrid Gin + Huma).
func NewRouter(d Deps) *gin.Engine {
	r, _ := NewRouterWithAPI(d)
	return r
}

// NewRouterWithAPI builds the engine and returns the Huma API (for OpenAPI tests).
func NewRouterWithAPI(d Deps) (*gin.Engine, huma.API) {
	if d.Config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	if d.States == nil {
		d.States = auth.NewStateStore(10*time.Minute, 10000)
	}
	if !d.DisableRateLimit {
		if d.GlobalLimiter == nil {
			d.GlobalLimiter = middleware.NewIPRateLimiter(60, 60)
		}
		if d.LoginLimiter == nil {
			d.LoginLimiter = middleware.NewIPRateLimiter(5, 5)
		}
		if d.StoreLimiter == nil {
			d.StoreLimiter = middleware.NewIPRateLimiter(10, 10)
		}
	}

	r := gin.New()
	r.Use(middleware.JSONRecovery(), gin.Logger(), middleware.CORS(d.Config.FrontendURL))
	if d.GlobalLimiter != nil {
		r.Use(d.GlobalLimiter.Middleware())
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Not Found"})
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"detail": "Method Not Allowed"})
	})

	apiGroup := r.Group("/api")
	apiGroup.GET("/health", health)

	// Pure Gin: OAuth redirects + plaintext storeapi (awkward in Huma).
	registerGinOnly(apiGroup, d)

	// Huma: JSON CRUD + OpenAPI docs (dev: /api/docs).
	api := humaapi.Register(r, apiGroup, humaapi.Deps{
		Config:      d.Config,
		Tokens:      d.Tokens,
		Users:       d.Users,
		Keys:        d.Keys,
		StoreTokens: d.StoreTokens,
		Profiles:    d.Profiles,
	})

	return r, api
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func registerGinOnly(g *gin.RouterGroup, d Deps) {
	authHandler := handlers.NewAuthHandler(d.Config, d.Tokens, d.States, d.GitHub, d.Auth)
	storeHandler := handlers.NewStoreAPIHandler(d.Config, d.StoreTokens)

	login := []gin.HandlerFunc{authHandler.GitHubLogin}
	if d.LoginLimiter != nil {
		login = append([]gin.HandlerFunc{d.LoginLimiter.Middleware()}, login...)
	}
	g.GET("/auth/github/login", login...)
	g.GET("/auth/github/callback", authHandler.GitHubCallback)
	g.POST("/auth/logout", authHandler.Logout)

	storeHandlers := []gin.HandlerFunc{storeHandler.Get}
	if d.StoreLimiter != nil {
		storeHandlers = append([]gin.HandlerFunc{d.StoreLimiter.Middleware()}, storeHandlers...)
	}
	g.GET("/storeapi", storeHandlers...)
}
