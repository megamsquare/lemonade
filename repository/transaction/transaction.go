package transaction

import (
	"context"
	"errors"
)

type Transaction struct {
	ID         int64   `json:"id"`
	SenderID   int64   `json:"sender_id"`
	ReceiverID int64   `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}

type NewTransaction struct {
	SenderID   int64   `json:"sender_id"`
	ReceiverID int64   `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}

type transactionMemoryStore struct {
	Table map[int64]*Transaction
}

var (
	ErrNotFound = errors.New("transactio not found")
)

type TransactionRepository interface {
	QueryByID(ctx context.Context, traceID string, id int64) (*Transaction, error)
	QueryAll(traceID string) ([]Transaction, error)
	Create(ctx context.Context, traceID string, newTransaction *NewTransaction) error
	Update(ctx context.Context, traceID string, user *Transaction) error
}

func NewTransactionMemoryStore() TransactionRepository {
	return &transactionMemoryStore{
		Table: make(map[int64]*Transaction),
	}
}
