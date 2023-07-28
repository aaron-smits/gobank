package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword compares a plaintext password to the encrypted password
func (a *Account) ComparePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

// Helper for making JWT token
// Creates a JWT token with the account number and user ID as claims
func makeJWTToken(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"account_number": account.AccountNumber,
		"user_id":        account.ID,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
		"token_type":     "Bearer",
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Helper for validating JWT token
// Parses the token and checks if it is valid based on the secret
func validateJWTToken(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

// Helper for getting user ID from JWT token claims
func getIDFromClaims(token *jwt.Token) (int, error) {
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("could not get user ID from claims")
	}
	return int(id), nil
}

// Middleware for JWT authentication
// 1. Validates the token
// 2. Checks if the user is an admin if the endpoint is admin-only
// 3. Checks if the user is accessing their own account by ID
// If any of the above checks fail, the middleware returns an error
func withJWTAuth(adminOnly bool, handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("call to JWT middleware")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "no token provided"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := validateJWTToken(tokenStr)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "error validating token"})
			return
		}
		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "token is invalid. unauthorized"})
			return
		}

		userID, err := getIDFromClaims(token)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "error getting user ID from token"})
			return
		}

		account, err := s.GetAccountByID(userID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "error getting account"})
			return
		}
		
		// Check if the user is an admin if the endpoint is admin-only
		if adminOnly && !account.IsAdmin {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "insufficient permissions"})
			return
		}

		// Check if the user is accessing their own account by ID
		if !account.IsAdmin && r.URL.Path != fmt.Sprintf("/account/%d", userID) {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "insufficient permissions"})
			return
		}

		handlerFunc(w, r)
	}
}
