package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func makeJWTToken(account *Account) (string, error) {
	claims := jwt.MapClaims{
		"account_number": account.AccountNumber,
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

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("call to JWT middleware")

		tokenStr := r.Header.Get("x-jwt-token")
		token, err := validateJWTToken(tokenStr)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}
		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}
		userID, err := getID(r)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token"})
			return
		}
		account, err := s.GetAccountByID(userID)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		accountNumber, ok := claims["account_number"].(float64)
		if !ok {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token"})
			return
		}
		account.AccountNumber = int64(accountNumber)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized token"})
			return
		}

		fmt.Println("claims", claims)
		handlerFunc(w, r)
	}
}
