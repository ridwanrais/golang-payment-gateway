package repository

import (
	"context"
	"time"

	"github.com/ridwanrais/golang-payment-gateway/internal/constants"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func (r *repository) InsertBrivaTransaction(ctx context.Context, data entity.BrivaData, referenceNumber, vaNumber string) (*entity.CreateBrivaResponse, error) {
	// Begin a database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Insert into the general transactions table
	var transactionID int
	var transactionUUID string
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (name, reference_number, transaction_date, transaction_amount, description, status, transaction_type)
         VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING transaction_id, uuid`, data.Nama,
		referenceNumber, time.Now(), data.Amount, data.Keterangan, constants.PAYMENT_PENDING, constants.PAYMENT_TYPE_VIRTUAL_ACCOUNT).Scan(&transactionID, &transactionUUID)
	if err != nil {
		return nil, err
	}

	// Insert into the virtual_account_transactions table and return its uuid
	var vaTransactionUUID string
	err = tx.QueryRow(ctx,
		`INSERT INTO virtual_account_transactions (transaction_id, bank_name, virtual_account_number, expiry_date)
         VALUES ($1, $2, $3, $4) RETURNING uuid`,
		transactionID, constants.BANK_NAME_BRI, vaNumber, data.ExpiredDate).Scan(&vaTransactionUUID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.CreateBrivaResponse{
		TransactionUUID:   transactionUUID,
		VaTransactionUUID: vaTransactionUUID,
	}, nil
}

