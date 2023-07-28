package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Reprsents the JSON API server
type APIServer struct {
	listenAddr string
	store      Storage
}

// NewAPIServer creates a new JSON API server
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the JSON API server and listens for requests
func (s *APIServer) Run() {
	// Create a new chi router and register the routes
	router := chi.NewRouter()
	// The account endpoint is for creating and getting accounts
	router.HandleFunc("/account", withJWTAuth(MakeHTTPHandlerFunc(s.handleAccount), s.store, true))
	// This endpoint is for getting and deleting accounts by ID
	router.HandleFunc("/account/{id}", withJWTAuth(MakeHTTPHandlerFunc(s.handleGetAccountByID), s.store, false))
	// This endpoint is for transferring money between accounts
	router.HandleFunc("/transfer", withJWTAuth(MakeHTTPHandlerFunc(s.handleTransfer), s.store, false))
	// This endpoint is for logging in and receiving a JWT token
	router.HandleFunc("/login", MakeHTTPHandlerFunc(s.handleLogin))
	log.Println("JSON API server running on", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// Log in and receive a JWT token
// Post to /login your account number and password
//
//	{
//		"account_number": 123456,
//		"password": "password"
//	}
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("unsupported method %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account, err := s.store.GetAccountByNumber(int(req.AccountNumber))
	if err != nil {
		return err
	}
	if !account.ComparePassword(req.Password) {
		return fmt.Errorf("unauthorized")
	}
	token, err := makeJWTToken(account)
	if err != nil {
		return err
	}
	resp := LoginResponse{
		AccountNumber: account.AccountNumber,
		Token:         token,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// Create a new account or get all accounts
// Post to /account to create a new account
//
//	{
//		"first_name": "John",
//		"last_name": "Doe",
//		"password": "password"
//		"is_admin": false
//	}
//
// Get /account to get all accounts
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid id %s", idStr)
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("unsupported method %s", r.Method)
}

// Get all accounts or get an account by ID
// This function is used in the handleAccount function to get all accounts when the endpoint
// is hit with the GET method
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// Create a new account
// Used in the handleAccount function to create a new account when the endpoint is hit with the POST method
// Example request body:
//
//	{
//		"first_name": "John",
//		"last_name": "Doe",
//		"password": "password"
//		"is_admin": false
//	}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	account, err := NewAccount(req.FirstName, req.LastName, req.Password, false)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %s", idStr)
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}
