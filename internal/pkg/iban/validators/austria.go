package validators

import (
	"unicode"
)

type Austria struct {
	baseValidator
}

func (a *Austria) ValidateLength(iban string) error {
	if len(iban) != 20 {
		return ErrIncorrectLength
	}

	return nil
}

func (a *Austria) ValidateBban(bban string) error {
	// Austrian BBAN must be numeric: https://www.xe.com/ibancalculator/sample/?ibancountry=austria
	for _, r := range bban {
		if !unicode.IsDigit(r) {
			return ErrIncorrectBBANFormat
		}
	}
	return nil
}
