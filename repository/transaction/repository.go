package transaction

import "context"

func(tr *transactionMemoryStore) QueryByID(ctx context.Context, traceID string, id int64) (*Transaction, error) {
	transaction, ok := tr.Table[id]
	if !ok {
		return nil, ErrNotFound
	}
	return transaction, nil
}

func(tr *transactionMemoryStore) QueryAll(traceID string) ([]Transaction, error) {
	values := []Transaction{}
	if len(tr.Table) < 1 {
		return nil, ErrNotFound
	}
	for _, value := range tr.Table {
		values = append(values, *value)
	}
	return values, nil
}

func(tr *transactionMemoryStore) Create(ctx context.Context, traceID string, newTransaction *NewTransaction) error {
	transaction := Transaction {
		ID: int64(len(tr.Table)),
		SenderID: newTransaction.SenderID,
		ReceiverID: newTransaction.ReceiverID,
		Amount: newTransaction.Amount,
	}
	tr.Table[transaction.ID] = &transaction

	return nil
}

func(tr *transactionMemoryStore) Update(ctx context.Context, traceID string, transaction *Transaction) error {
	savedTransaction, err := tr.QueryByID(ctx, traceID, transaction.ID)
	if err != nil {
		return err
	}
	tr.Table[savedTransaction.ID] = savedTransaction

	return nil
}