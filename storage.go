package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Storage is an interface for storing and retrieving accounts
// All of these methods are required to be implemented
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

// PostgresStore is an implementation of the Storage interface
// The only thing in this struct is a pointer to a sql.DB
// This is for connecting to the database
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgresStore and connects to the database
func NewPostgresStore() (*PostgresStore, error) {
	connectionString := os.Getenv("POSTGRES_URL")
	if connectionString == "" {
		log.Fatal("POSTGRES_URL environment variable not set")
	}
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

// CreateAccountTable creates the accounts table if it doesn't exist\
// This is called in main.go when the server starts
func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE if not exists accounts(
		id SERIAL PRIMARY KEY,
		first_name varchar(50) NOT NULL,
		last_name varchar(50) NOT NULL,
		account_number BIGINT NOT NULL,
		encrypted_password varchar(100) NOT NULL,
		balance BIGINT NOT NULL,
		created_at timestamp
		)`
	_, err := s.db.Exec(query)
	if err == nil {
		fmt.Println("Account Table Initialized")
	}
	return err
}

// CreateAccount creates a new account in the database
// Takes a pointer to an account
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO accounts (
			first_name,
			last_name,
			account_number,
			encrypted_password,
			balance,
			created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6
			)`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.AccountNumber,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	defer resp.Close()
	return nil
}

// DeleteAccount deletes an account from the database
// In the future, this should probably be a soft delete
// I should create a new column called deleted_at and set it to the current time
// Or I could create a new table called deleted_accounts and move the account there
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAccount updates an account in the database
// This is not implemented yet
// Ideas for implementation:
// parameterize this function so that it can take a map of fields to update
// or take a pointer to an account and update all of the fields
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

// GetAccountByID gets an account from the database by ID
// This is used in the handleAccountByID function in api.go
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	account := new(Account)
	query := `SELECT * FROM accounts WHERE id=$1`
	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.AccountNumber,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByNumber gets an account from the database by account number
// This is used in the handleAccount function in api.go
func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE account_number=$1", number)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.AccountNumber,
			&account.EncryptedPassword,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		return account, nil
	}

	return nil, fmt.Errorf("account not found")
}

// This gets all of the accounts from the database
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.AccountNumber,
			&account.EncryptedPassword,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
