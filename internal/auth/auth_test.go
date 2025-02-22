package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"net/http"
)

func TestValidateJWT(t *testing.T) {
	tokenSecret := "secret"
	userID := uuid.New()
	expiresIn := time.Hour

	// Create a valid token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Test valid token
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedUserID)
	}

	// Test invalid token
	invalidToken := token + "invalid"
	_, err = ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Error("Expected error for invalid token, got none")
	}

	// Test expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}
	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Error("Expected error for expired token, got none")
	}
}

func TestGetBearerToken(t *testing.T) {
	// Test valid token
	tokenSecret := "secret"
	userID := uuid.New()
	expiresIn := time.Hour
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
	bearerToken := http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	parsedToken, err := GetBearerToken(bearerToken)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if parsedToken != token {
		t.Errorf("Expected token %v, got %v", token, parsedToken)
	}
}
