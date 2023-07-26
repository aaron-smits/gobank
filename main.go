package main

import (
	"fmt"
	"log"
	"flag"
)

func seedAccount(store Storage, firstName, lastName, password string) *Account{
	account, err := NewAccount(firstName, lastName, password)
	if err != nil {
		log.Fatal(err)
	}
	
	if err := store.CreateAccount(account); err != nil {
		log.Fatal(err)
	}

	return account
}

func seedAccounts(s Storage) {
	seedAccount(s, "John", "Doe", "password")
}

func main() {
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
