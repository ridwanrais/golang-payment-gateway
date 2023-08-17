package repository

import (
	"context"
	"time"

	"github.com/ridwanrais/golang-payment-gateway/internal/constants"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func (r *repository) InsertBrivaTransaction(ctx context.Context, data entity.BrivaData, referenceNumber, vaNumber string) (int, error) {
	// Begin a database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Insert into the general transactions table
	var transactionID int
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (reference_number, transaction_date, transaction_amount, description, status, transaction_type)
         VALUES ($1, $2, $3, $4, $5, $6) RETURNING transaction_id`,
		referenceNumber, time.Now(), data.Amount, data.Keterangan, constants.PAYMENT_PENDING, constants.PAYMENT_TYPE_VIRTUAL_ACCOUNT).Scan(&transactionID)
	if err != nil {
		return 0, err
	}

	// Insert into the virtual_account_transactions table
	_, err = tx.Exec(ctx,
		`INSERT INTO virtual_account_transactions (transaction_id, bank_name, virtual_account_number, expiry_date)
         VALUES ($1, $2, $3, $4)`,
		transactionID, constants.BANK_NAME_BRI, vaNumber, data.ExpiredDate)
	if err != nil {
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}
