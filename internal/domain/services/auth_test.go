package services

import (
	"os"
	"testing"
	"time"

	"github.com/cHeLoVe4uK/EM_Project/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGetToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("SECRET_KEY")

	userID, username := "123", "testuser"
	token, err := GetToken(userID, username)
	require.NoError(t, err, "Error should not occur while generating token")
	require.NotEmpty(t, token, "Token should not be empty")

	claims := &models.Claims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("test_secret_key"), nil
	})
	require.NoError(t, err, "Token parsing should not fail")
	require.Equal(t, userID, claims.UserID, "UserID should match")
	require.Equal(t, username, claims.Username, "Username should match")
}

func TestAuthenticate(t *testing.T) {
	os.Setenv("SECRET_KEY", "test_secret_key")
	defer os.Unsetenv("SECRET_KEY")

	validClaims := models.Claims{
		UserID:   "123",
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Chat App",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := GetToken(validClaims.UserID, validClaims.Username)
	require.NoError(t, err, "Error should no occur while generating token")

	claims, err := Authenticate(token)
	require.NoError(t, err, "No error should occur while authenticating")
	require.Equal(t, validClaims.UserID, claims.UserID, "UserID should match")
	require.Equal(t, validClaims.Username, claims.Username, "Username should match")

	invalidToken := "Invalid token"
	_, err = Authenticate(invalidToken)
	require.Error(t, err, "Error should occur with invalid token")
}
