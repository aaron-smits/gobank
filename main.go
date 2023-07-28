package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func seedAccount(s Storage, firstName, lastName, password string, isAdmin bool) *Account {
	account, err := NewAccount(firstName, lastName, password, isAdmin)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.CreateAccount(account); err != nil {
		log.Fatal(err)
	}
	if !isAdmin {
		fmt.Println("NEW ACCOUNT SEEDED:", account.AccountNumber)
	}
	if isAdmin {
		fmt.Println("NEW ADMIN ACCOUNT SEEDED:", account.AccountNumber)
	}
	return account
}

func seedAccounts(s Storage) {
	seedAccount(s, "John", "Doe", "password", false)
	seedAccount(s, "Defacto", "Admin", "password", true)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()
	store, err := NewPostgresStore()
	fmt.Println("Postgres Store Created")
	if err != nil {
		log.Fatal(err)
	}
	if err = store.CreateAccountTable(); err != nil {
		log.Fatal(err)
	}
	if *seed {
		fmt.Println("Seeding Database")
		seedAccounts(store)
	}
	server := NewAPIServer(":5555", store)
	server.Run()
	fmt.Println("Server Running on :5555")
}
