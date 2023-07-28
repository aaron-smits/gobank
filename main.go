package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func seedAccount(s Storage, firstName, lastName, password string, isAdmin bool, balance ...int64) *Account {
	var bal int64
	if len(balance) > 0 {
		bal = balance[0]
	}
	account, err := NewAccount(firstName, lastName, password, isAdmin, bal)
	if err != nil {
		log.Fatal(err)
	}

	acc, err := s.CreateAccount(account); 
	if err != nil {
		log.Fatal(err)
	}
	if !isAdmin {
		fmt.Println("NEW ACCOUNT SEEDED:", account.AccountNumber)
	}
	if isAdmin {
		fmt.Println("NEW ADMIN ACCOUNT SEEDED:", account.AccountNumber)
	}
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "John", "Doe", "password", false, 1000)
	seedAccount(s, "Cool", "Guy", "password", false, 1000)
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
	fmt.Println("Account Table Created")
	if *seed {
		fmt.Println("Seeding Database")
		seedAccounts(store)
	}

	server := NewAPIServer(":5555", store)
	server.Run()
	fmt.Println("Server Running on :5555")
}
