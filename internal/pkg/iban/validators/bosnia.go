package validators

import (
	"unicode"
)

type Bosnia struct {
	baseValidator
}

func (a *Bosnia) ValidateLength(iban string) error {
	if len(iban) != 20 {
		return ErrIncorrectLength
	}

	return nil
}

func (a *Bosnia) ValidateBban(bban string) error {
	// Bosnian BBAN must be numeric: https://www.xe.com/ibancalculator/sample/?ibancountry=austria
	for _, r := range bban {
		if !unicode.IsDigit(r) {
			return ErrIncorrectBBANFormat
		}
	}

	// XXX: Bosnian BBAN supposedly uses a non-variant of ISO7064 MOD-97-10 to validate the BBAN check digits
	//      however, the check-digits cannot be computed successfully, neither using the entire BBAN,
	//      nor excluding the check-digits, nor using only the account number,
	//      both with the standard-mod97 and with the 98-r complement.
	//      I cannot find any reliable information in english on the exact specs
	//      for Bosnian BBAN checksum computation, so I'll leave it for now.

	// Bosnian BBAN uses ISO 7064 MOD-97-10 as a checksum algorithm
	//if 98-mod97(bban[:14]) != 1 {
	//	return ErrIncorrectBBANChecksum
	//}

	return nil
}
