package validators

import (
	"regexp"
)

var UKBbanRegex = regexp.MustCompile(`^[A-Z]{4}\d{14}$`)

type GB struct {
	baseValidator
}

func (a *GB) ValidateLength(iban string) error {
	if len(iban) != 22 {
		return ErrIncorrectLength
	}

	return nil
}

func (a *GB) ValidateBban(bban string) error {
	// GB BBAN must follow this specific format: https://www.xe.com/ibancalculator/sample/?ibancountry=united-kingdom
	if !UKBbanRegex.MatchString(bban) {
		return ErrIncorrectBBANFormat
	}
	return nil
}
