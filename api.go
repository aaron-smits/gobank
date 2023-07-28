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
	// The account endpoint is for creating and getting accounts. Admins only
	router.HandleFunc("/accounts", withJWTAuth(true, MakeHTTPHandlerFunc(s.handleAccounts), s.store))
	// This endpoint is for getting and deleting accounts by ID
	router.HandleFunc("/account/{id}", withJWTAuth(false, MakeHTTPHandlerFunc(s.handleGetAccountByID), s.store))
	// This endpoint is for transferring money between accounts. Admins only.
	router.HandleFunc("/transfer", withJWTAuth(true, MakeHTTPHandlerFunc(s.handleTransfer), s.store))
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
		return fmt.Errorf("invalid request body")
	}

	account, err := s.store.GetAccountByNumber(int(req.AccountNumber))
	if err != nil {
		return fmt.Errorf("unauthorized")
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
func (s *APIServer) handleAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
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

	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}

	return fmt.Errorf("unsupported method %s", r.Method)
}

// Get all accounts or get an account by ID
// This function is used in the handleAccount function to get all accounts when the endpoint
// is hit with the GET method
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return fmt.Errorf("error getting accounts: %v", err)
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
		return fmt.Errorf("invalid request body")
	}
	account, err := NewAccount(
		req.FirstName, 
		req.LastName, 
		req.Password,  
		false, 
		req.Balance,
	)
	if err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}
	if err := s.store.CreateAccount(account); err != nil {
		return fmt.Errorf("error creating account: %v", err)
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
		return fmt.Errorf("error deleting account: %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// Shape of the request body for updating an account
//{
//	"first_name": "John",
//	"last_name": "Doe",
//	"account_number": 123456,
//	"is_admin": false
//}
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %s", idStr)
	}

	req := new(UpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("invalid request body")
	}
	acc, err := s.store.GetAccountByID(id)
	if err != nil {
		return fmt.Errorf("error getting account: %v", err)
	}

	updatedAccount := &Account{
		ID:                acc.ID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		AccountNumber:     req.AccountNumber,
		IsAdmin: 		   req.IsAdmin,
	}

	_, err = s.store.UpdateAccountByID(id, updatedAccount)
	if err != nil {
		return fmt.Errorf("error updating account: %v", err)
	}

	return WriteJSON(w, http.StatusOK, updatedAccount)


}

//Request body sample
//{
//	"to_account_id": 123456,
//	"from_account_id": 123456,
//	"amount": 100
//}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return fmt.Errorf("invalid request body")
	}


	acc, err := s.store.MakeTransfer(
		transferReq.ToAccountID, 
		transferReq.FromAccountID, 
		transferReq.Amount,
	)
	
	if err != nil {
		return fmt.Errorf("error making transfer: %v", err)
	}


	return WriteJSON(w, http.StatusOK, fmt.Sprintf("Transfer successful. New balance: %d", acc.Balance))
}
