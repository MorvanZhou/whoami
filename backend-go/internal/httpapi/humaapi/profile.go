package humaapi

import (
	"context"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/danielgtaylor/huma/v2"

	"whoami/backend-go/internal/httpapi/middleware"
	"whoami/backend-go/internal/service"
)

type profileBodyInput struct {
	Body struct {
		Content string `json:"content" maxLength:"5000" doc:"Markdown profile content"`
	}
}

type profileOutput struct {
	Body struct {
		Content *string `json:"content"`
	}
}

type messageOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

type historyItem struct {
	Timestamp    string `json:"timestamp"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
}

type historyOutput struct {
	Body []historyItem
}

type rollbackInput struct {
	Timestamp string `path:"timestamp"`
}

func registerProfiles(
	api huma.API,
	d Deps,
	jwtMW func(huma.Context, func(huma.Context)),
	apiKeyMW func(huma.Context, func(huma.Context)),
) {
	// Agent-facing (/profile) — API key
	huma.Register(api, huma.Operation{
		OperationID: "get-profile",
		Method:      http.MethodGet,
		Path:        "/profile",
		Summary:     "Get profile (API key)",
		Tags:        []string{"profile"},
		Middlewares: huma.Middlewares{apiKeyMW},
		Security:    []map[string][]string{{"apiKeyAuth": {}}},
	}, getProfile(d))

	huma.Register(api, huma.Operation{
		OperationID: "update-profile",
		Method:      http.MethodPost,
		Path:        "/profile",
		Summary:     "Update profile (API key)",
		Tags:        []string{"profile"},
		Middlewares: huma.Middlewares{apiKeyMW},
		Security:    []map[string][]string{{"apiKeyAuth": {}}},
	}, updateProfile(d))

	huma.Register(api, huma.Operation{
		OperationID: "list-profile-history",
		Method:      http.MethodGet,
		Path:        "/profile/history",
		Summary:     "List profile history",
		Tags:        []string{"profile"},
		Middlewares: huma.Middlewares{apiKeyMW},
		Security:    []map[string][]string{{"apiKeyAuth": {}}},
	}, func(ctx context.Context, _ *struct{}) (*historyOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		items, err := d.Profiles.ListHistory(ctx, user.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to list history")
		}
		out := &historyOutput{Body: make([]historyItem, 0, len(items))}
		for _, it := range items {
			out.Body = append(out.Body, historyItem{
				Timestamp:    it.Timestamp,
				Size:         it.Size,
				LastModified: it.LastModified,
			})
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "rollback-profile",
		Method:      http.MethodPost,
		Path:        "/profile/rollback/{timestamp}",
		Summary:     "Rollback profile to a history version",
		Tags:        []string{"profile"},
		Middlewares: huma.Middlewares{apiKeyMW},
		Security:    []map[string][]string{{"apiKeyAuth": {}}},
	}, func(ctx context.Context, in *rollbackInput) (*messageOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		ok, err := d.Profiles.Rollback(ctx, user.ID, in.Timestamp)
		if err != nil {
			if strings.Contains(err.Error(), "invalid timestamp") {
				return nil, huma.Error400BadRequest("Invalid timestamp")
			}
			return nil, huma.Error500InternalServerError("Failed to rollback profile")
		}
		if !ok {
			return nil, huma.Error404NotFound("History version not found")
		}
		out := &messageOutput{}
		out.Body.Message = "Profile rolled back to " + in.Timestamp
		return out, nil
	})

	// Console-facing (/web/profile) — JWT
	huma.Register(api, huma.Operation{
		OperationID: "get-web-profile",
		Method:      http.MethodGet,
		Path:        "/web/profile",
		Summary:     "Get profile (JWT)",
		Tags:        []string{"web-profile"},
		Middlewares: huma.Middlewares{jwtMW},
		Security:    []map[string][]string{{"bearerAuth": {}}},
	}, getProfile(d))

	huma.Register(api, huma.Operation{
		OperationID: "update-web-profile",
		Method:      http.MethodPost,
		Path:        "/web/profile",
		Summary:     "Update profile (JWT)",
		Tags:        []string{"web-profile"},
		Middlewares: huma.Middlewares{jwtMW},
		Security:    []map[string][]string{{"bearerAuth": {}}},
	}, updateProfile(d))
}

func getProfile(d Deps) func(context.Context, *struct{}) (*profileOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*profileOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		content, err := d.Profiles.Get(ctx, user.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to get profile")
		}
		out := &profileOutput{}
		out.Body.Content = content
		return out, nil
	}
}

func updateProfile(d Deps) func(context.Context, *profileBodyInput) (*messageOutput, error) {
	return func(ctx context.Context, in *profileBodyInput) (*messageOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		if utf8.RuneCountInString(in.Body.Content) > service.ProfileMaxLength {
			return nil, huma.Error400BadRequest("content exceeds max length")
		}
		if err := d.Profiles.Save(ctx, user.ID, in.Body.Content); err != nil {
			if strings.Contains(err.Error(), "exceeds") {
				return nil, huma.Error400BadRequest(err.Error())
			}
			return nil, huma.Error500InternalServerError("Failed to update profile")
		}
		out := &messageOutput{}
		out.Body.Message = "Profile updated"
		return out, nil
	}
}
