package entity

import (
	"database/sql"
	"time"
)

type Transaction struct {
	TransactionID     int32          `db:"transaction_id"`
	UUID              string         `db:"uuid"`
	ReferenceNumber   string         `db:"reference_number"`
	TransactionDate   time.Time      `db:"transaction_date"`
	TransactionAmount float64        `db:"transaction_amount"`
	Description       sql.NullString `db:"description"`
	Status            string         `db:"status"`
	TransactionType   string         `db:"transaction_type"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}

type VirtualAccountTransaction struct {
	VATransactionID      int32          `db:"va_transaction_id"`
	UUID                 string         `db:"uuid"`
	TransactionID        int32          `db:"transaction_id"`
	BankName             string         `db:"bank_name"`
	VirtualAccountNumber string         `db:"virtual_account_number"`
	ExpiryDate           sql.NullTime   `db:"expiry_date"`
	Metadata             sql.NullString `db:"metadata"` // Assuming this will be serialized to/from string
}

type GetVirtualAccountResponse struct {
	ReferenceNumber      string     `json:"referenceNumber"`
	TransactionDate      time.Time  `json:"transactionDate"`
	TransactionAmount    int        `json:"transactionAmount"`
	Description          string     `json:"description"`
	PaymentStatus        string     `json:"paymentStatus"`
	BankName             string     `json:"bankName"`
	VirtualAccountNumber string     `json:"virtualAccountNumber"`
	ExpiryDate           *time.Time `json:"expiryDate"`
	Metadata             string    `json:"metadata"`
}
