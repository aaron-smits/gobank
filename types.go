package main

import "math/rand"

type Account struct {
	ID 	int
	Name 	string
	Number 	int64
	Balance int64
}

func NewAccount(name string) *Account {
	return &Account{
		ID:      rand.Intn(1000),
		Name:    name,
		Number:  int64(rand.Intn(10000000)),
	}
}