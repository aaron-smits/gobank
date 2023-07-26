package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TransferRequest struct {
	ToAccountID int `json:"to_account_id"`
	Amount    int `json:"amount"`
}

type Account struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	EncryptedPassword string    `json:"-"`
	AccountNumber int64     `json:"account_number"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		EncryptedPassword: string(encpw),
		AccountNumber: int64(rand.Intn(1000000)),
		Balance:       0,
		CreatedAt:     time.Now().UTC(),
	}, nil
	
	
}

type LoginRequest struct {
	AccountNumber int64 `json:"account_number"`
	Password      string `json:"password"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}
