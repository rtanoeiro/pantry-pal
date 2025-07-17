package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var TestConfig = Config{
	Secret: "SuperTestSecret",
}

func TestHashing(t *testing.T) {
	passwordMap := []string{
		"admin",
		"mysecurepassword",
	}

	for _, password := range passwordMap {
		_, errorHash := HashPassword(password)
		if errorHash != nil {
			t.Errorf("Expected password %s to be hashable, but got an error instead.", password)
		}
	}
}

func TestCompareHashing(t *testing.T) {
	passwordMap := map[string]string{
		"admin":  "$2a$10$2WP4ssk27deQ7hGKWhwYl.DUlN740Gc0jDQwUT7eHIR8qXUcKnCw2",
		"123456": "$2a$10$ICmo9SYxEmZgZNKAfzTzi.xmZDlqAeqsPesBNTb6kRsBpxLvD22tm",
	}
	for password, hashedPassword := range passwordMap {
		errorHash := CheckPasswordHash(password, hashedPassword)
		if errorHash != nil {
			t.Errorf("Expected password %s to comparable to %s, but got an error instead.", password, hashedPassword)
		}
	}
}

func TestValidateDate(t *testing.T) {
	dateMap := map[string]bool{
		"2000-01-01": false,
		"2999-12-12": true,
	}
	for date, expectedResults := range dateMap {
		result := ValidateDate(date)
		if result != expectedResults {
			t.Errorf("Date %s failed to evaluate. Expected %t, got %t.", date, expectedResults, result)
		}
	}
}

func TestJWTTokens(t *testing.T) {
	expireTime := 60 * time.Second
	userMap := []string{
		"123456",
		"123456-78901-23456",
	}
	for _, user := range userMap {
		token, errorToken := MakeJWT(user, TestConfig.Secret, expireTime)
		if errorToken != nil {
			t.Errorf("Failed to generate JWT for user %s: %v", user, errorToken)
			continue
		}
		request := httptest.NewRequest(http.MethodGet, "/login", nil)
		request.AddCookie(&http.Cookie{
			Name:     "JWTToken",
			Value:    token,
			Expires:  time.Now().Add(expireTime),
			HttpOnly: true,
		})

		userID, errorUser := GetUserIDFromTokenAndValidate(request, &TestConfig)

		if errorUser != nil {
			t.Errorf("Failed to get user %s from JWT. Error: %s", user, errorUser)
		}
		if userID != user {
			t.Errorf("Got different user than expected. Expected %s, got %s", user, userID)
			continue
		}
	}
}
