package entity

import (
	"database/sql"
	"time"
)

type Transaction struct {
	TransactionID     int32          `db:"transaction_id"`
	UUID              string         `db:"uuid"`
	ReferenceNumber   string         `db:"reference_number"`
	Name              sql.NullString `db:"name"`
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
	Name                 string     `json:"name"`
	TransactionDate      time.Time  `json:"transactionDate"`
	TransactionAmount    int        `json:"transactionAmount"`
	Description          string     `json:"description"`
	PaymentStatus        string     `json:"paymentStatus"`
	BankName             string     `json:"bankName"`
	VirtualAccountNumber string     `json:"virtualAccountNumber"`
	ExpiryDate           *time.Time `json:"expiryDate"`
	Metadata             string     `json:"metadata"`
}

type UpdateVaRequest struct {
	TransactionUUID   string `json:"-"`
	VaTransactionUUID string `json:"vaUuid"`
	PhoneNumber       string `json:"phoneNumber"`
	Name              string `json:"name"`
	Amount            int    `json:"amount"`
	Note              string `json:"note"`
	// PaymentStatus        string `json:"paymentStatus"`
	BankName             string `json:"bankName"`
	VirtualAccountNumber string `json:"-"`
	ExpiryDate           string `json:"expiryDate"`
}

type UpdateVaResponse struct {
	ReferenceNumber      string `json:"referenceNumber"`
	VirtualAccountNumber string `json:"virtualAccountNumber"`
	TransactionUUID      string `json:"transactionUuid"`
	VaTransactionUUID    string `json:"vaTransactionUuid"`
}
