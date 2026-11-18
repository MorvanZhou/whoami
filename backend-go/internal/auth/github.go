package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubUser is the normalized profile from GitHub APIs.
type GitHubUser struct {
	ID        string
	Email     *string
	Name      *string
	AvatarURL *string
}

// GitHubOAuth handles the GitHub authorization-code flow.
type GitHubOAuth struct {
	config     *oauth2.Config
	httpClient *http.Client
	userURL   string
	emailsURL  string
}

func NewGitHubOAuth(clientID, clientSecret, redirectURI string) *GitHubOAuth {
	return &GitHubOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
		httpClient: http.DefaultClient,
		userURL:    "https://api.github.com/user",
		emailsURL:  "https://api.github.com/user/emails",
	}
}

// SetHTTPClient overrides the HTTP client (useful in tests).
func (g *GitHubOAuth) SetHTTPClient(c *http.Client) {
	if c != nil {
		g.httpClient = c
	}
}

// SetEndpoints overrides OAuth and API endpoints (useful in tests).
func (g *GitHubOAuth) SetEndpoints(authURL, tokenURL, userURL, emailsURL string) {
	g.config.Endpoint = oauth2.Endpoint{AuthURL: authURL, TokenURL: tokenURL}
	if userURL != "" {
		g.userURL = userURL
	}
	if emailsURL != "" {
		g.emailsURL = emailsURL
	}
}

func (g *GitHubOAuth) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (g *GitHubOAuth) ExchangeAndFetchUser(ctx context.Context, code string) (*GitHubUser, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, g.httpClient)
	tok, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("github token exchange: %w", err)
	}

	userReq, err := http.NewRequestWithContext(ctx, http.MethodGet, g.userURL, nil)
	if err != nil {
		return nil, err
	}
	userReq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	userReq.Header.Set("Accept", "application/json")

	userResp, err := g.httpClient.Do(userReq)
	if err != nil {
		return nil, fmt.Errorf("github user: %w", err)
	}
	defer userResp.Body.Close()
	if userResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(userResp.Body, 2048))
		return nil, fmt.Errorf("github user status %d: %s", userResp.StatusCode, body)
	}

	var raw struct {
		ID        any    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	id := stringifyID(raw.ID)
	if id == "" {
		return nil, fmt.Errorf("github user missing id")
	}

	email := strings.TrimSpace(raw.Email)
	if email == "" {
		email, err = g.fetchPrimaryEmail(ctx, tok.AccessToken)
		if err != nil {
			return nil, err
		}
	}

	name := strings.TrimSpace(raw.Name)
	if name == "" {
		name = strings.TrimSpace(raw.Login)
	}

	out := &GitHubUser{ID: id}
	if email != "" {
		out.Email = &email
	}
	if name != "" {
		out.Name = &name
	}
	if raw.AvatarURL != "" {
		av := raw.AvatarURL
		out.AvatarURL = &av
	}
	return out, nil
}

func (g *GitHubOAuth) fetchPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.emailsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Email != "" {
			return e.Email, nil
		}
	}
	return "", nil
}

func stringifyID(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case json.Number:
		return t.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
