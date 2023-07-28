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
	UpdateAccountByID(int, *Account) (*Account, error)
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
	GetAdminStatus(int) (bool, error)
	MakeTransfer(int, int, int) (*Account, error)
	AddBalanceTx(*sql.Tx, int, int) error
	SubtractBalanceTx(*sql.Tx, int, int) error
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
		created_at timestamp,
		is_admin boolean DEFAULT false
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
			created_at,
			is_admin
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7
			)`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.AccountNumber,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
		acc.IsAdmin,
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
	// Check that the account exists
	_, err := s.GetAccountByID(id)
	if err != nil {
		return err
	}
	
	// Delete the account
	query := `DELETE FROM accounts WHERE id=$1`
	_, err = s.db.Query(query, id)
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
func (s *PostgresStore) UpdateAccountByID(id int, accountDetails *Account) (*Account, error) {
	// Make sure the account exists
	query := `
		UPDATE accounts
		SET 
		first_name=$1, 
		last_name=$2, 
		account_number=$3,
		is_admin=$4
		WHERE id=$5`
	//Check to make sure accountDetails is not nil
	if accountDetails == nil {
		return nil, fmt.Errorf("account details cannot be nil")
	}
	
	_, err := s.db.Query(
		query,
		accountDetails.FirstName,
		accountDetails.LastName,
		accountDetails.AccountNumber,
		accountDetails.IsAdmin,
		id,
	)

	if err != nil {
		return nil, err
	}

	// Get the updated account from the database
	account, err := s.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	return account, nil
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
		&account.IsAdmin,
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
			&account.IsAdmin,
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
			&account.IsAdmin,
		)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// GetAdminStatus gets the admin status of an account
// This is used in the withAdminAuth middleware in auth.go
func (s *PostgresStore) GetAdminStatus(id int) (bool, error) {
	var isAdmin bool
	row := s.db.QueryRow("SELECT is_admin FROM accounts WHERE id=$1", id)
	err := row.Scan(&isAdmin)
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}

func (s *PostgresStore) GetBalanceTx(tx *sql.Tx, id int) (int64, error) {
    var balance int64
    row := tx.QueryRow("SELECT balance FROM accounts WHERE id=$1", id)
    err := row.Scan(&balance)
    if err != nil {
        return 0, err
    }

    return balance, nil
}

func (s *PostgresStore) AddBalanceTx(tx *sql.Tx, id int, amount int) error {
    balance, err := s.GetBalanceTx(tx, id)
    if err != nil {
        return err
    }

    balance += int64(amount)
    _, err = tx.Exec("UPDATE accounts SET balance=$1 WHERE id=$2", balance, id)
    if err != nil {
        return err
    }

    return nil
}

func (s *PostgresStore) SubtractBalanceTx(tx *sql.Tx, id int, amount int) error {
    balance, err := s.GetBalanceTx(tx, id)
    if err != nil {
        return err
    }

    balance -= int64(amount)
    _, err = tx.Exec("UPDATE accounts SET balance=$1 WHERE id=$2", balance, id)
    if err != nil {
        return err
    }

    return nil
}

// MakeTransfer makes a transfer from one account to another
// checks the balance of the from account and subtracts the amount from their balance
// then adds the amount to the to account
// This is used in the handleTransfer function in api.go
func (s *PostgresStore) MakeTransfer(toAcc, fromAcc, amount int) (*Account, error) {
	// Begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	// Get the balance of the from account
	fromBalance, err := s.GetBalanceTx(tx, fromAcc)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Check if the from account has enough money
	if fromBalance < int64(amount) {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient funds")
	}

	// Subtract the amount from the from account
	err = s.SubtractBalanceTx(tx, fromAcc, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}


	// Add the amount to the to account
	err = s.AddBalanceTx(tx, toAcc, amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Get the updated account from the database
	account, err := s.GetAccountByID(fromAcc)
	if err != nil {
		return nil, err
	}

	return account, nil
}
