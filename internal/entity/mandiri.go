package entity

type Amount struct {
	Value    string `json:"value"` // Assuming that the value is a string in the format "number.number"
	Currency string `json:"currency"`
}

type BillDetail struct {
	BillAmount Amount `json:"billAmount"`
}

type MandiriVaData struct {
	PartnerServiceId    string       `json:"partnerServiceId"`
	CustomerNo          string       `json:"customerNo"`
	VirtualAccountNo    string       `json:"virtualAccountNo"`
	VirtualAccountName  string       `json:"virtualAccountName"`
	VirtualAccountEmail string       `json:"virtualAccountEmail,omitempty"` // omitempty for optional fields
	VirtualAccountPhone string       `json:"virtualAccountPhone,omitempty"` // omitempty for optional fields
	TrxId               string       `json:"trxId"`
	TotalAmount         Amount       `json:"totalAmount"`
	BillDetails         []BillDetail `json:"billDetails"`
	ExpiredDate         string       `json:"expiredDate"`
}

type MandiriResponseData struct {
	ResponseCode       string        `json:"responseCode"`
	ResponseMessage    string        `json:"responseMessage,omitempty"`
	VirtualAccountData MandiriVaData `json:"virtualAccountData"`
}

type MandiriAccessTokenResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	AccessToken     string `json:"accessToken"`
	TokenType       string `json:"tokenType"`
	ExpiresIn       string `json:"expiresIn"`
}
