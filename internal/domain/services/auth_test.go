package services

import (
	"context"
	"os"
	"testing"

	"github.com/cHeLoVe4uK/EM_Project/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGetTokens(t *testing.T) {
	os.Setenv("SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("SECRET_KEY")

	s := NewService()
	user := models.User{
		UserID:   "123",
		Username: "testuser",
	}

	tokens, err := s.GetTokens(context.Background(), user)
	require.NoError(t, err, "Error should not occur while generating token")
	require.NotEmpty(t, tokens.AccessToken, "Token should not be empty")

	claims := &models.Claims{}
	_, err = jwt.ParseWithClaims(tokens.AccessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("test_secret_key"), nil
	})
	require.NoError(t, err, "Token parsing should not fail")
	require.Equal(t, user.UserID, claims.UserID, "UserID should match")
	require.Equal(t, user.Username, claims.Username, "Username should match")
}

func TestAuthenticate(t *testing.T) {
	os.Setenv("SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("SECRET_KEY")

	s := NewService()
	user := models.User{
		UserID:   "123",
		Username: "testuser",
	}

	tokens, err := s.GetTokens(context.Background(), user)
	require.NoError(t, err, "Error should no occur while generating token")

	claims, err := s.Authenticate(context.Background(), tokens)
	require.NoError(t, err, "No error should occur while authenticating")
	require.Equal(t, user.UserID, claims.UserID, "UserID should match")
	require.Equal(t, user.Username, claims.Username, "Username should match")

	_, err = s.Authenticate(context.Background(), models.Tokens{AccessToken: "Invalid token"})
	require.Error(t, err, "Error should occur with invalid token")
}
