package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword compares a plaintext password to the encrypted password
func (a *Account) ComparePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func makeJWTToken(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"account_number": account.AccountNumber,
		"user_id":        account.ID,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWTToken(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

// GetIDFromClaims is a helper function for getting the ID from JWT token claims
func getIDFromClaims(token *jwt.Token) (int, error) {
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("could not get user ID from claims")
	}
	return int(id), nil
}

func getAccountNumberFromClaims(token *jwt.Token) (int64, error) {
	claims := token.Claims.(jwt.MapClaims)
	accountNumber, ok := claims["account_number"].(float64)
	if !ok {
		return 0, fmt.Errorf("could not get account number from claims")
	}
	return int64(accountNumber), nil
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage, adminOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("call to JWT middleware")

		tokenStr := r.Header.Get("x-jwt-token")
		token, err := validateJWTToken(tokenStr)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "error validating token"})
			return
		}
		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "token is invalid. unauthorized"})
			return
		}

		userID, err := getIDFromClaims(token)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token. Could not get user ID"})
			return
		}
		account, err := s.GetAccountByID(userID)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token"})
			return
		}

		// Check if the user is an admin if the endpoint is admin-only
		if adminOnly && !account.IsAdmin {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}

		// Check if the user is accessing their own account by ID
		if !account.IsAdmin && r.URL.Path != fmt.Sprintf("/account/%d", userID) {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}

		account.AccountNumber, err = getAccountNumberFromClaims(token)

		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token"})
			return
		}
		fmt.Printf("account number: %d\n", account.AccountNumber)
		fmt.Printf("user id: %d\n", account.ID)
		handlerFunc(w, r)
	}
}
