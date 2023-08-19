package repository

import (
	"context"

	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func (r *repository) GetVaTransaction(ctx context.Context, virtualAccountUuuid string) (*entity.Transaction, *entity.VirtualAccountTransaction, error) {
	const query = `
	SELECT
		t.transaction_id, t.uuid, t.reference_number, t.transaction_date, t.transaction_amount, t.description, t.status, t.transaction_type, t.created_at, t.updated_at,
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

	if err := r.db.QueryRow(ctx, query, virtualAccountUuuid).Scan(
		&tx.TransactionID, &tx.UUID, &tx.ReferenceNumber, &tx.TransactionDate, &tx.TransactionAmount, &tx.Description, &tx.Status, &tx.TransactionType, &tx.CreatedAt, &tx.UpdatedAt,
		&vaTx.VATransactionID, &vaTx.UUID, &vaTx.TransactionID, &vaTx.BankName, &vaTx.VirtualAccountNumber, &vaTx.ExpiryDate, &vaTx.Metadata,
	); err != nil {
		return nil, nil, err
	}

	return tx, vaTx, nil
}
