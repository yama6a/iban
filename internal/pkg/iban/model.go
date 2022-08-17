package iban

// swagger:model
type IBAN struct {
	CountryCode string `json:"country_code"`
	CheckDigits uint   `json:"check_digit"`
	BBAN        string `json:"bban"`
}

// swagger:model
type httpResponse struct {
	Error   *string `json:"error"`
	IsValid bool    `json:"is_valid"`
	IBAN    *IBAN   `json:"iban"`
}
