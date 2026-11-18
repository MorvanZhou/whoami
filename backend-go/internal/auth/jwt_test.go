package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestTokenManager_CreateParseTable(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		expireDays int
		userID     string
		mutate     func(token string) string
		wantErr    bool
	}{
		{name: "roundtrip", secret: "s1", expireDays: 7, userID: "user-1"},
		{name: "default expire when zero", secret: "s1", expireDays: 0, userID: "user-2"},
		{
			name: "tampered token", secret: "s1", expireDays: 1, userID: "user-3",
			mutate: func(token string) string {
				return token + "x"
			},
			wantErr: true,
		},
		{
			name: "wrong secret", secret: "s1", expireDays: 1, userID: "user-4",
			mutate: func(token string) string {
				other := NewTokenManager("other", 1)
				tok, err := other.Create("user-4")
				if err != nil {
					t.Fatal(err)
				}
				return tok
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTokenManager(tt.secret, tt.expireDays)
			tok, err := m.Create(tt.userID)
			if err != nil {
				t.Fatal(err)
			}
			if tt.mutate != nil {
				tok = tt.mutate(tok)
			}
			got, err := m.ParseUserID(tok)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil || got != tt.userID {
				t.Fatalf("got=%q err=%v", got, err)
			}
		})
	}
}

func TestTokenManager_Expired(t *testing.T) {
	m := NewTokenManager("secret", 7)
	claims := jwt.MapClaims{
		"sub": "u1",
		"exp": time.Now().UTC().Add(-time.Hour).Unix(),
		"iat": time.Now().UTC().Add(-2 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = m.ParseUserID(signed)
	if err == nil {
		t.Fatal("expected expired error")
	}
	if !strings.Contains(err.Error(), "invalid token") {
		t.Fatalf("err=%v", err)
	}
}

func TestTokenManager_MaxAgeSeconds(t *testing.T) {
	if NewTokenManager("s", 3).MaxAgeSeconds() != 3*86400 {
		t.Fatal("unexpected max age")
	}
}
