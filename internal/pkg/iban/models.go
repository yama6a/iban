package iban

import "fmt"

// swagger:model
type IBAN struct {
	CountryCode string `json:"country_code"`
	CheckDigits string `json:"check_digits"`
	BBAN        string `json:"bban"`
}

func (i IBAN) String() string {
	return fmt.Sprintf("%s%s%s", i.CountryCode, i.CheckDigits, i.BBAN)
}

// swagger:model
type httpResponse struct {
	Error   *string `json:"error"`
	IsValid bool    `json:"is_valid"`
	IBAN    *IBAN   `json:"iban"`
}
