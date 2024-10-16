package db

import (
	"database/sql"
)

type Store struct {
	*Queries
	db *sql.DB
}

type TransferTransactionParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Ammount       int64 `json:"ammount"`
}

type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var transactionKey = struct{}{}
