package humaapi

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"whoami/backend-go/internal/httpapi/middleware"
	"whoami/backend-go/internal/service"
)

type createKeyInput struct {
	Body struct {
		Label string `json:"label" minLength:"1" maxLength:"15" doc:"Human-readable key label"`
	}
}

type createKeyOutput struct {
	Status int
	Body   struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		StoreURL  string `json:"store_url"`
		ExpiresIn int    `json:"expires_in"`
	}
}

type apiKeyItem struct {
	ID        string  `json:"id"`
	KeyPrefix string  `json:"key_prefix"`
	KeySuffix string  `json:"key_suffix"`
	Label     *string `json:"label"`
	CreatedAt string  `json:"created_at"`
	LastUsed  *string `json:"last_used"`
}

type listKeysOutput struct {
	Body []apiKeyItem
}

type keyIDInput struct {
	KeyID string `path:"key_id"`
}

type keyStatusOutput struct {
	Body struct {
		Configured bool `json:"configured"`
	}
}

type deleteKeyOutput struct {
	Status int
}

func registerKeys(api huma.API, d Deps, jwtMW func(huma.Context, func(huma.Context))) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-key",
		Method:        http.MethodPost,
		Path:          "/keys",
		Summary:       "Create API key",
		Tags:          []string{"keys"},
		DefaultStatus: http.StatusCreated,
		Middlewares:   huma.Middlewares{jwtMW},
		Security:      []map[string][]string{{"bearerAuth": {}}},
	}, func(ctx context.Context, in *createKeyInput) (*createKeyOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		plain, key, err := d.Keys.Create(ctx, user.ID, in.Body.Label)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to create API key")
		}
		token, err := d.StoreTokens.Create(ctx, user.ID, plain, key.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to create store token")
		}
		label := ""
		if key.Label != nil {
			label = *key.Label
		}
		out := &createKeyOutput{Status: http.StatusCreated}
		out.Body.ID = key.ID
		out.Body.Label = label
		out.Body.StoreURL = d.Config.FrontendURL + "/api/storeapi?token=" + token
		out.Body.ExpiresIn = service.StoreTokenExpireMinutes * 60
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-keys",
		Method:      http.MethodGet,
		Path:        "/keys",
		Summary:     "List API keys",
		Tags:        []string{"keys"},
		Middlewares: huma.Middlewares{jwtMW},
		Security:    []map[string][]string{{"bearerAuth": {}}},
	}, func(ctx context.Context, _ *struct{}) (*listKeysOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		keys, err := d.Keys.List(ctx, user.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to list API keys")
		}
		out := &listKeysOutput{Body: make([]apiKeyItem, 0, len(keys))}
		for _, k := range keys {
			item := apiKeyItem{
				ID:        k.ID,
				KeyPrefix: k.KeyPrefix,
				KeySuffix: k.KeySuffix,
				Label:     k.Label,
			}
			if !k.CreatedAt.IsZero() {
				item.CreatedAt = k.CreatedAt.UTC().Format(time.RFC3339Nano)
			}
			if k.LastUsed != nil {
				s := k.LastUsed.UTC().Format(time.RFC3339Nano)
				item.LastUsed = &s
			}
			out.Body = append(out.Body, item)
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-key",
		Method:        http.MethodDelete,
		Path:          "/keys/{key_id}",
		Summary:       "Delete API key",
		Tags:          []string{"keys"},
		DefaultStatus: http.StatusNoContent,
		Middlewares:   huma.Middlewares{jwtMW},
		Security:      []map[string][]string{{"bearerAuth": {}}},
	}, func(ctx context.Context, in *keyIDInput) (*deleteKeyOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		ok, err := d.Keys.Revoke(ctx, in.KeyID, user.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to delete API key")
		}
		if !ok {
			return nil, huma.Error404NotFound("API key not found")
		}
		return &deleteKeyOutput{Status: http.StatusNoContent}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-key-status",
		Method:      http.MethodGet,
		Path:        "/keys/{key_id}/status",
		Summary:     "Check whether store token was consumed",
		Tags:        []string{"keys"},
		Middlewares: huma.Middlewares{jwtMW},
		Security:    []map[string][]string{{"bearerAuth": {}}},
	}, func(ctx context.Context, in *keyIDInput) (*keyStatusOutput, error) {
		user := middleware.UserFromContext(ctx)
		if user == nil {
			return nil, huma.Error401Unauthorized("Not authenticated")
		}
		consumed, err := d.StoreTokens.IsConsumed(ctx, in.KeyID, user.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to check status")
		}
		out := &keyStatusOutput{}
		out.Body.Configured = consumed
		return out, nil
	})
}
