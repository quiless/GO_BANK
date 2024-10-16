package db

import (
	"context"
	"database/sql"
	"fmt"
)

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("Transaction error: %v, Rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *Store) TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error) {
	var result TransferTransactionResult

	err := store.execTransaction(ctx, func(q *Queries) error {

		var err error

		transactionName := ctx.Value(transactionKey)

		fmt.Println(transactionName, "creating transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Ammount,
		})

		if err != nil {
			return err
		}

		fmt.Println(transactionName, "creating entry 1")

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Ammount,
		})

		if err != nil {
			return err
		}

		fmt.Println(transactionName, "creating entry 2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    +arg.Ammount,
		})

		if err != nil {
			return err
		}

		// get account -> update values

		fmt.Println(transactionName, "getting account 1 for update")

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Ammount,
		})
		if err != nil {
			return err
		}

		fmt.Println(transactionName, "updating account 2")

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Ammount,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
