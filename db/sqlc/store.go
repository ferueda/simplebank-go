package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
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

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) DeleteAccountTx(ctx context.Context, accountId int64) error {
	err := s.execTrx(ctx, func(q *Queries) error {
		var err error

		err = q.DeleteEntry(ctx, accountId)
		if err != nil {
			return err
		}

		err = q.DeleteTransfer(ctx, accountId)
		if err != nil {
			return err
		}

		err = q.DeleteAccount(ctx, accountId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTrx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
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
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

func addMoney(ctx context.Context, q *Queries, fromAccId, toAccId, fromAmount, toAmount int64) (fromAcc, toAcc Account, err error) {
	fromAcc, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccId,
		Amount: fromAmount,
	})
	if err != nil {
		return
	}

	toAcc, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccId,
		Amount: toAmount,
	})
	if err != nil {
		return
	}
	return
}

func (s *Store) execTrx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
