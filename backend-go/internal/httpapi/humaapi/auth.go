package humaapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"whoami/backend-go/internal/httpapi/middleware"
)

type meOutput struct {
	Body struct {
		ID        string  `json:"id"`
		Email     *string `json:"email"`
		Name      *string `json:"name"`
		AvatarURL *string `json:"avatar_url"`
		Provider  string  `json:"provider"`
		CreatedAt *string `json:"created_at"`
	}
}

func registerAuthMe(api huma.API, jwtMW func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID: "get-auth-me",
		Method:      http.MethodGet,
		Path:        "/auth/me",
		Summary:     "Get current user",
		Tags:        []string{"auth"},
		Middlewares: huma.Middlewares{jwtMW},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, _ *struct{}) (*meOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		out := &meOutput{}
		out.Body.ID = user.ID
		out.Body.Email = user.Email
		out.Body.Name = user.Name
		out.Body.AvatarURL = user.AvatarURL
		out.Body.Provider = user.Provider
		if !user.CreatedAt.IsZero() {
			s := user.CreatedAt.UTC().Format("2006-01-02T15:04:05.999999Z07:00")
			out.Body.CreatedAt = &s
		}
		return out, nil
	})
}
