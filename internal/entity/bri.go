package entity

type BrivaData struct {
	InstitutionCode string  `json:"institutionCode"`
	BrivaNo         string  `json:"brivaNo"`
	CustCode        string  `json:"custCode"`
	Nama            string  `json:"nama"`
	Amount          string  `json:"amount"`
	Keterangan      string  `json:"keterangan"`
	StatusBayar     string  `json:"statusBayar"`
	ExpiredDate     string  `json:"expiredDate"`
	LastUpdate      *string `json:"lastUpdate,omitempty"`
}

type BriResponseData struct {
	Status              bool      `json:"status"`
	ResponseCode        string    `json:"responseCode"`
	ResponseDescription string    `json:"responseDescription,omitempty"`
	ErrDesc             string    `json:"errDesc,omitempty"`
	Data                BrivaData `json:"data"`
}

type CreateBrivaRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Note        string `json:"note"`
}

type CreateBrivaResponse struct {
	ReferenceNumber      string `json:"referenceNumber"`
	VirtualAccountNumber string `json:"virtualAccountNumber"`
	// TransactionID        int32  `json:"transactionId"`
	TransactionUUID string `json:"transactionUuid"`
	// VaTransactionID      int32  `json:"vaTransactionId"`
	VaTransactionUUID string `json:"vaTransactionUuid"`
}
