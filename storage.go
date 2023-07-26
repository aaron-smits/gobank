package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE if not exists accounts(
		id SERIAL PRIMARY KEY,
		first_name varchar(50) NOT NULL,
		last_name varchar(50) NOT NULL,
		account_number BIGINT NOT NULL,
		balance BIGINT NOT NULL,
		created_at timestamp
		)`
	_, err := s.db.Exec(query)
	if err == nil {
		fmt.Println("Account Table Initialized")
	}
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO accounts (
			first_name,
			last_name,
			account_number,
			balance,
			created_at
			) VALUES (
				$1, $2, $3, $4, $5
				)`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.AccountNumber,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	defer resp.Close()
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	account := new(Account)
	query := `SELECT * FROM accounts WHERE id=$1`
	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.AccountNumber,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error){
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
