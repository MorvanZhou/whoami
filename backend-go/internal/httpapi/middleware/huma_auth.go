package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"whoami/backend-go/internal/auth"
	"whoami/backend-go/internal/domain"
	"whoami/backend-go/internal/service"
)

// HumaRequireJWT authenticates via access_token cookie or Authorization Bearer JWT.
func HumaRequireJWT(api huma.API, tokens *auth.TokenManager, users domain.UserRepository) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		tokenStr := ""
		if c, err := huma.ReadCookie(ctx, "access_token"); err == nil && c != nil && c.Value != "" {
			tokenStr = c.Value
		}
		if tokenStr == "" {
			authz := ctx.Header("Authorization")
			if strings.HasPrefix(authz, "Bearer ") {
				tokenStr = strings.TrimSpace(authz[7:])
			}
		}
		if tokenStr == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not authenticated")
			return
		}

		userID, err := tokens.ParseUserID(tokenStr)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid token")
			return
		}

		user, err := users.FindByID(ctx.Context(), userID)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal error")
			return
		}
		if user == nil || !user.IsActive {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "User not found")
			return
		}

		next(huma.WithValue(ctx, ContextUserKey, user))
	}
}

// HumaRequireAPIKey authenticates via Authorization: Bearer <api_key>.
func HumaRequireAPIKey(api huma.API, keys *service.KeyService) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authz := ctx.Header("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Missing API key")
			return
		}
		plain := strings.TrimSpace(authz[7:])
		user, err := keys.Authenticate(ctx.Context(), plain)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal error")
			return
		}
		if user == nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid API key")
			return
		}
		next(huma.WithValue(ctx, ContextUserKey, user))
	}
}

// UserFromContext returns the authenticated user set by Huma auth middleware.
func UserFromContext(ctx context.Context) *domain.User {
	u, _ := ctx.Value(ContextUserKey).(*domain.User)
	return u
}
