package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccountID int `json:"to_account_id"`
	Amount    int `json:"amount"`
}

type Account struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AccountNumber int64     `json:"account_number"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		AccountNumber: int64(rand.Intn(10000000)),
		CreatedAt:     time.Now().UTC(),
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
