package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute SQL queries and transactions
type Store struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Convert arg to CreateTransferParams
		createTransferParams := CreateTransferParams(arg)

		result.Transfer, err = q.CreateTransfer(ctx, createTransferParams)
		if err != nil {
			return err
		}

		createFromEntryParams := CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}

		result.FromEntry, err = q.CreateEntry(ctx, createFromEntryParams)

		if err != nil {
			return err
		}

		createToEntryParams := CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}

		result.ToEntry, err = q.CreateEntry(ctx, createToEntryParams)

		if err != nil {
			return err
		}

		// TODO: update accounts balance

		return nil
	})

	return result, err
}
