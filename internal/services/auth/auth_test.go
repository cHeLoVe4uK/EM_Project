package services

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/cHeLoVe4uK/EM_Project/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGetTokens(t *testing.T) {
	os.Setenv("TOKEN_SALT", "test_secret_key")
	os.Setenv("ACCESS_TOKEN_EXP", "24")
	defer os.Unsetenv("TOKEN_SALT")
	defer os.Unsetenv("ACCESS_TOKEN_EXP")

	ate, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	require.NoError(t, err, "No error should occuer while converting ACCESS_TOKEN_EXP to int")
	s := NewService(os.Getenv("TOKEN_SALT"), ate)
	user := models.User{
		UserID:   "123",
		Username: "testuser",
	}

	tokens, err := s.GetTokens(context.Background(), user)
	require.NoError(t, err, "Error should not occur while generating token")
	require.NotEmpty(t, tokens.AccessToken, "Token should not be empty")

	claims := &Claims{}
	_, err = jwt.ParseWithClaims(tokens.AccessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("test_secret_key"), nil
	})
	require.NoError(t, err, "Token parsing should not fail")
	require.Equal(t, user.UserID, claims.UserID, "UserID should match")
	require.Equal(t, user.Username, claims.Username, "Username should match")
}

func TestAuthenticate(t *testing.T) {
	os.Setenv("TOKEN_SALT", "test_secret_key")
	os.Setenv("ACCESS_TOKEN_EXP", "24")
	defer os.Unsetenv("TOKEN_SALT")
	defer os.Unsetenv("ACCESS_TOKEN_EXP")

	ate, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	require.NoError(t, err, "No error should occuer while converting ACCESS_TOKEN_EXP to int")
	s := NewService(os.Getenv("TOKEN_SALT"), ate)
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
