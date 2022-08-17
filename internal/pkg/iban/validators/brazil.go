package validators

import (
	"regexp"
)

var BrazilBbanRegex = regexp.MustCompile(`^\d{23}[A-Z]{1}[A-Z\d]{1}$`)

type Brazil struct {
	baseValidator
}

func (a *Brazil) ValidateLength(iban string) error {
	if len(iban) != 29 {
		return ErrIncorrectLength
	}

	return nil
}

func (a *Brazil) ValidateBban(bban string) error {
	// Brazilian BBAN must follow this specific format: https://www.xe.com/ibancalculator/sample/?ibancountry=brazil
	if !BrazilBbanRegex.MatchString(bban) {
		return ErrIncorrectBBANFormat
	}
	return nil
}
