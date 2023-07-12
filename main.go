package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgresStore()
	fmt.Println("Postgres Store Created")
	if err != nil {
		log.Fatal(err)
	}
	if err = store.CreateAccountTable(); err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":5555", store)
	server.Run()
	fmt.Println("Server Running on :5555")
}
