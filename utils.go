package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// this type represents a function that takes a request and returns an error
type apiFunc func(w http.ResponseWriter, r *http.Request) error

// this function is for the router to be able take requests and return responses
func MakeHTTPHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// This is a helper function for writing JSON responses
// It takes a status code and a value and writes the JSON response
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// NewAccount creates a new account and hashes the password
// This function is used in the seedAccounts function in main.go
// Currently the account number is a random number between 0 and 1,000,000
// In the future we will want to make sure that the account number is unique
func NewAccount(firstName, lastName, password string, isAdmin bool,  balance ...int64) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var bal int64
	if len(balance) > 0 {
		bal = balance[0]
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encpw),
		AccountNumber:     int64(rand.Intn(1000000)),
		Balance:           bal,
		CreatedAt:         time.Now().UTC(),
		IsAdmin:           isAdmin,
	}, nil

}
