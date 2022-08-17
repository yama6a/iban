package iban

import (
	"fmt"
	"regexp"
)

var (
	_ Parser = &Service{}

	ibanRegexp                 = regexp.MustCompile(`^([A-Z]{2})(\d{2})([A-Z\d]+)$`)
	ErrIncorrectIbanFormat     = fmt.Errorf("provided string does not satisfy the iban format: %s", ibanRegexp.String())
	ErrCountryCodeNotSupported = fmt.Errorf("country code is not supported")
	ErrBBANEmpty               = fmt.Errorf("BBAN is empty")
	ErrCountryCodeEmpty        = fmt.Errorf("country code is empty")
)

type Service struct {
	validators map[string]countryValidator
}

func NewService() *Service {
	return &Service{
		validators: countryValidators,
	}
}

func (svc *Service) Parse(ibanStr string) (IBAN, error) {
	matches := ibanRegexp.FindStringSubmatch(ibanStr)
	if matches == nil || len(matches) != 4 {
		return IBAN{}, ErrIncorrectIbanFormat
	}

	countryCode, checkDigits, bban := matches[1], matches[2], matches[3]

	return IBAN{
		CountryCode: countryCode,
		CheckDigits: checkDigits,
		BBAN:        bban,
	}, nil
}

// Validate validates the iban's format and checks the check-digits.
func (svc *Service) Validate(i IBAN) error {
	if i.CountryCode == "" {
		return ErrCountryCodeEmpty
	}
	if i.BBAN == "" {
		return ErrBBANEmpty
	}

	validator, ok := svc.validators[i.CountryCode]
	if !ok {
		return ErrCountryCodeNotSupported
	}

	err := validator.ValidateIbanLength(i)
	if err != nil {
		return fmt.Errorf("iban length validation error: %w", err)
	}

	err = validator.ValidateIbanChecksum(i)
	if err != nil {
		return fmt.Errorf("iban checksum validation error: %w", err)
	}

	err = validator.ValidateBbanFormat(i)
	if err != nil {
		return fmt.Errorf("bban format validation error: %w", err)
	}

	err = validator.ValidateBbanChecksum(i)
	if err != nil {
		return fmt.Errorf("bban checksum validation error: %w", err)
	}

	return nil
}
