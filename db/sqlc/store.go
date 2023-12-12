package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccount: arg.FromAccountID,
			ToAccount:   arg.ToAccountID,
			Amount:      arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = plusMoney(ctx, q, -arg.Amount, arg.FromAccountID, arg.Amount, arg.ToAccountID)
		} else {
			result.ToAccount, result.FromAccount, err = plusMoney(ctx, q, arg.Amount, arg.ToAccountID, -arg.Amount, arg.FromAccountID)
		}
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func plusMoney(
	ctx context.Context,
	q *Queries,
	amount1 int64,
	accountID1 int64,
	amount2 int64,
	accountID2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAmountToAccount(ctx, AddAmountToAccountParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAmountToAccount(ctx, AddAmountToAccountParams{
		Amount: amount2,
		ID:     accountID2,
	})

	return
}
