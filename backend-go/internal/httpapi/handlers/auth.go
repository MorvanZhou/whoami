package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/config"
	"whoami/backend-go/internal/service"
)

var allowedRedirects = map[string]struct{}{
	"/dashboard": {},
	"/settings":  {},
	"/keys":      {},
}

type AuthHandler struct {
	cfg      *config.Config
	tokens   *auth.TokenManager
	states   *auth.StateStore
	github   *auth.GitHubOAuth
	authSvc  *service.AuthService
}

func NewAuthHandler(
	cfg *config.Config,
	tokens *auth.TokenManager,
	states *auth.StateStore,
	github *auth.GitHubOAuth,
	authSvc *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		cfg:     cfg,
		tokens:  tokens,
		states:  states,
		github:  github,
		authSvc: authSvc,
	}
}

func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	if h.cfg.GitHubClientID == "" || h.cfg.GitHubClientSecret == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"detail": "GitHub OAuth is not configured"})
		return
	}

	redirect := c.Query("redirect")
	state, err := h.states.Put(redirect)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create OAuth state"})
		return
	}
	c.Redirect(http.StatusFound, h.github.AuthCodeURL(state))
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Missing code or state"})
		return
	}

	redirectTarget, ok := h.states.Take(state)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid state"})
		return
	}

	ghUser, err := h.github.ExchangeAndFetchUser(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"detail": "GitHub OAuth failed"})
		return
	}

	user, err := h.authSvc.FindOrCreateUser(
		c.Request.Context(),
		"github",
		ghUser.ID,
		ghUser.Email,
		ghUser.Name,
		ghUser.AvatarURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create user"})
		return
	}

	jwt, err := h.tokens.Create(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create session"})
		return
	}

	h.redirectWithCookie(c, jwt, redirectTarget)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", h.cfg.CookieSecure, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (h *AuthHandler) redirectWithCookie(c *gin.Context, jwtToken, redirectTarget string) {
	if _, ok := allowedRedirects[redirectTarget]; !ok {
		redirectTarget = ""
	}
	redirectURL := h.cfg.FrontendURL + "/auth/callback"
	if redirectTarget != "" {
		redirectURL += "?redirect=" + url.QueryEscape(redirectTarget)
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"access_token",
		jwtToken,
		h.tokens.MaxAgeSeconds(),
		"/",
		"",
		h.cfg.CookieSecure,
		true, // HttpOnly
	)
	c.Redirect(http.StatusFound, redirectURL)
}
