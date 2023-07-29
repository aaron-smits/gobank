package main

import (
	"net/http"
	"time"
)

// this type represents a function that takes a request and returns an error
type apiFunc func(w http.ResponseWriter, r *http.Request) error

// Api error type for server error responses
type ApiError struct {
	Error string `json:"error"`
}

// TransferRequest is the request body for the transfer endpoint
type TransferRequest struct {
	ToAccountID   int `json:"to_account_id"`
	FromAccountID int `json:"from_account_id"`
	Amount        int `json:"amount"`
}

// Account is the model for storing account information
type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	EncryptedPassword string    `json:"-"`
	AccountNumber     int64     `json:"account_number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
	IsAdmin           bool      `json:"is_admin"`
}

type LoginRequest struct {
	AccountNumber int64  `json:"account_number"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	AccountNumber int64  `json:"sub"` // Sub is part of the JWT spec
	Token         string `json:"access_token"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Balance   int64  `json:"balance"`
	IsAdmin   bool   `json:"is_admin"`
}

type UpdateAccountRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AccountNumber int64  `json:"account_number"`
	IsAdmin       bool   `json:"is_admin"`
}
