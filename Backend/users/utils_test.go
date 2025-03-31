package main

import (
	"os"
	"testing"
	"time"
	"users/models"
	"users/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT tests that we can generate a valid JWT from a user struct
func TestGenerateJWT(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "testsecret")

	user := models.User{
		UserID: 42,
		Name:   "Jason Bourne",
		Email:  "jason@ufl.edu",
		Phone:  "5551994567",
	}

	tokenString, err := utils.GenerateJWT(user)
	assert.NoError(t, err, "should generate JWT without error")
	assert.NotEmpty(t, tokenString, "tokenString should not be empty")
}

// TestParseJWTValidToken checks that a valid token will parse correctly
func TestParseJWTValidToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "testsecret")

	user := models.User{
		UserID: 42,
		Name:   "Jason Bourne",
		Email:  "jason@ufl.edu",
		Phone:  "5551994567",
	}

	// Generate a token
	tokenString, err := utils.GenerateJWT(user)
	assert.NoError(t, err)

	token, parseErr := utils.ParseJWT(tokenString)
	assert.NoError(t, parseErr)
	assert.True(t, token.Valid, "Expected token to be valid")

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Claims should be of type jwt.MapClaims")
	userMap, ok := claims["user"].(map[string]interface{})
	assert.True(t, ok, "Expected 'user' key in claims")

	assert.Equal(t, float64(42), userMap["UserID"], "UserID should match")
	assert.Equal(t, "Jason Bourne", userMap["Name"], "Name should match")
	assert.Equal(t, "jason@ufl.edu", userMap["Email"], "Email should match")
}

// TestParseJWTInvalidToken checks that an invalid token fails
func TestParseJWTInvalidToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "testsecret")
	invalidToken := "abc.def.ghi"

	token, err := utils.ParseJWT(invalidToken)
	assert.Error(t, err, "Should fail to parse invalid token")
	assert.Nil(t, token, "Token should be nil on parse error")
}

func TestExpiredToken(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "testsecret")

	claims := jwt.MapClaims{
		"user": map[string]interface{}{
			"UserID": 123,
			"Name":   "Expired User",
			"Email":  "expired@ufl.edu",
		},
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("testsecret"))
	assert.NoError(t, err, "Signing a test token should not fail")

	parsedToken, parseErr := utils.ParseJWT(tokenString)
	assert.Nil(t, parsedToken, "Parsed token should be nil for an expired token")
	assert.Error(t, parseErr, "Should raise an error for an expired token")
}
