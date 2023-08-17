package entity

type Response struct {
	Status       bool                   `json:"status"`
	ResponseCode int                    `json:"responseCode"`
	Message      string                 `json:"message"`
	Data         map[string]interface{} `json:"data"`
}
