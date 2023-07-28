package main

import "time"

// Api error type for server error responses
type ApiError struct {
	Error string `json:"error"`
}

// TransferRequest is the request body for the transfer endpoint
type TransferRequest struct {
	ToAccountID int `json:"to_account_id"`
	Amount      int `json:"amount"`
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
	AccountNumber int64  `json:"account_number"`
	Token         string `json:"token"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"is_admin"`
}
