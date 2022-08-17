package iban

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

var (
	ErrIncorrectLength       = errors.New("IBAN has the incorrect length for the specified country")
	ErrIncorrectBBANFormat   = errors.New("IBAN has the incorrect BBAN format for the specified country")
	ErrIncorrectBBANChecksum = errors.New("IBAN has the incorrect BBAN checksum for the specified country")
	ErrIncorrectIBANChecksum = errors.New("IBAN has the incorrect checksum")
)

var (
	countryValidators = map[string]countryValidator{
		"AL": {CountryCode: "AL", Length: 28, BBANRegex: regexp.MustCompile(`^[0-9]{8}[0-9A-Z]{16}$`)},
		"AT": {CountryCode: "AT", Length: 20, BBANRegex: regexp.MustCompile(`^\d{16}$`)},
		"BA": {
			CountryCode: "BA",
			Length:      20,
			BBANRegex:   regexp.MustCompile(`^\d{16}$`),
			BBANChecksumFunc: func(bban string) bool {
				return true

				// return 98-mod97(bban[:14]) != 1
				// XXX: Bosnian BBAN supposedly uses a non-variant of ISO7064 MOD-97-10 to validate the BBAN check digits
				//      however, the check-digits cannot be computed successfully, neither using the entire BBAN,
				//      nor excluding the check-digits, nor using only the account number,
				//      both with the standard-mod97 and with the 98-r complement.
				//      I cannot find any reliable information in english on the exact specs
				//      for Bosnian BBAN checksum computation, so I'll leave it for now.

			},
		},
		"BR": {CountryCode: "BR", Length: 29, BBANRegex: regexp.MustCompile(`^\d{23}[A-Z]{1}[A-Z\d]{1}$`)},
		"CH": {CountryCode: "CH", Length: 21, BBANRegex: regexp.MustCompile(`^\d{17}$`)},
		"DE": {CountryCode: "DE", Length: 22, BBANRegex: regexp.MustCompile(`^\d{18}$`)},
		"FR": {CountryCode: "FR", Length: 27, BBANRegex: regexp.MustCompile(`^\d{10}[A-Z0-9]{11}\d{2}$`)},
		"GB": {CountryCode: "GB", Length: 22, BBANRegex: regexp.MustCompile(`^[A-Z]{4}\d{14}$`)},
	}
)

type countryValidator struct {
	CountryCode      string
	Length           int
	BBANRegex        *regexp.Regexp
	BBANChecksumFunc func(string) bool
}

func (c countryValidator) ValidateIbanLength(iban IBAN) error {
	if len(iban.String()) != c.Length {
		return ErrIncorrectLength
	}

	return nil
}

func (c countryValidator) ValidateBbanFormat(iban IBAN) error {
	if !c.BBANRegex.MatchString(iban.BBAN) {
		return ErrIncorrectBBANFormat
	}

	return nil
}

func (c countryValidator) ValidateBbanChecksum(iban IBAN) error {
	if c.BBANChecksumFunc != nil && !c.BBANChecksumFunc(iban.BBAN) {
		return ErrIncorrectBBANChecksum
	}

	return nil
}

func (c countryValidator) ValidateIbanChecksum(iban IBAN) error {
	// Ref: https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN
	transposedIban := fmt.Sprintf("%s%s%s", iban.BBAN, iban.CountryCode, iban.CheckDigits)

	str := ""
	for _, r := range transposedIban {
		if unicode.IsDigit(r) {
			str += string(r)
			continue
		}

		digits, err := convertLetterToInt(r)
		if err != nil {
			return fmt.Errorf("failed to convert IBAN into numeric format: %w", err)
		}
		str += strconv.Itoa(digits)
	}

	if mod97(str) != 1 {
		return ErrIncorrectIBANChecksum
	}

	return nil
}
