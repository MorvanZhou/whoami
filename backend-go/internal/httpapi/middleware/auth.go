package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/service"
)

const ContextUserKey = "user"

// RequireJWT authenticates via access_token cookie or Authorization: Bearer <jwt>.
func RequireJWT(tokens *auth.TokenManager, users domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		if cookie, err := c.Cookie("access_token"); err == nil {
			tokenStr = cookie
		}
		if tokenStr == "" {
			authz := c.GetHeader("Authorization")
			if strings.HasPrefix(authz, "Bearer ") {
				tokenStr = strings.TrimSpace(authz[7:])
			}
		}
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
			return
		}

		userID, err := tokens.ParseUserID(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token"})
			return
		}

		user, err := users.FindByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Internal error"})
			return
		}
		if user == nil || !user.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "User not found"})
			return
		}

		c.Set(ContextUserKey, user)
		c.Next()
	}
}

// RequireAPIKey authenticates via Authorization: Bearer <api_key>.
func RequireAPIKey(keys *service.KeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Missing API key"})
			return
		}
		plain := strings.TrimSpace(authz[7:])
		user, err := keys.Authenticate(c.Request.Context(), plain)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Internal error"})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Invalid API key"})
			return
		}
		c.Set(ContextUserKey, user)
		c.Next()
	}
}

// CurrentUser extracts the authenticated user from Gin context.
func CurrentUser(c *gin.Context) *domain.User {
	v, ok := c.Get(ContextUserKey)
	if !ok {
		return nil
	}
	user, _ := v.(*domain.User)
	return user
}
