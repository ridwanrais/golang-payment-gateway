package repository

import (
	"context"

	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func (r *repository) GetVaTransaction(ctx context.Context, virtualAccountUuid string) (*entity.Transaction, *entity.VirtualAccountTransaction, error) {
	const query = `
	SELECT
		t.transaction_id, t.uuid, t.name, t.reference_number, t.transaction_date, t.transaction_amount, t.description, t.status, t.transaction_type, t.created_at, t.updated_at,
		va.va_transaction_id, va.uuid, va.transaction_id, va.bank_name, va.virtual_account_number, va.expiry_date, va.metadata
	FROM
		transactions t
	JOIN
		virtual_account_transactions va ON t.transaction_id = va.transaction_id
	WHERE
		va.uuid = $1
	`

	tx := &entity.Transaction{}
	vaTx := &entity.VirtualAccountTransaction{}

	if err := r.db.QueryRow(ctx, query, virtualAccountUuid).Scan(
		&tx.TransactionID, &tx.UUID, &tx.Name, &tx.ReferenceNumber, &tx.TransactionDate, &tx.TransactionAmount, &tx.Description, &tx.Status, &tx.TransactionType, &tx.CreatedAt, &tx.UpdatedAt,
		&vaTx.VATransactionID, &vaTx.UUID, &vaTx.TransactionID, &vaTx.BankName, &vaTx.VirtualAccountNumber, &vaTx.ExpiryDate, &vaTx.Metadata,
	); err != nil {
		return nil, nil, err
	}

	return tx, vaTx, nil
}

func (r *repository) UpdateVaTransaction(ctx context.Context, updateData entity.UpdateVaRequest) (*entity.UpdateVaResponse, error) {
	// Begin a database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Update the general transactions table
	var transactionID int
	var transactionUUID string
	err = tx.QueryRow(ctx,
		`UPDATE transactions SET name=$1, transaction_amount=$2, description=$3
	 WHERE uuid=$4 RETURNING transaction_id, uuid`,
		updateData.Name, updateData.Amount, updateData.Note, updateData.TransactionUUID).Scan(&transactionID, &transactionUUID)
	if err != nil {
		return nil, err
	}

	// Update the virtual_account_transactions table
	var vaTransactionUUID string
	err = tx.QueryRow(ctx,
		`UPDATE virtual_account_transactions SET bank_name=$1, virtual_account_number=$2, expiry_date=$3
		 WHERE transaction_id=$4 RETURNING uuid`,
		updateData.BankName, updateData.VirtualAccountNumber, updateData.ExpiryDate, transactionID).Scan(&vaTransactionUUID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.UpdateVaResponse{
		TransactionUUID:   transactionUUID,
		VaTransactionUUID: vaTransactionUUID,
	}, nil
}
