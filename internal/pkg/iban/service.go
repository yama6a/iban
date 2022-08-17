package iban

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/ymakhloufi/pfc/internal/pkg/iban/validators"
)

var (
	_ Parser = &Service{}

	ibanRegexp = regexp.MustCompile(`^*([A-Z]{2})(\d{2})([A-Z\d]+)$`)
)

type CountryValidator interface {
	ValidateChecksum(string) error
	ValidateLength(string) error
	ValidateBban(string) error
}

type Service struct {
	validators map[string]CountryValidator
}

func NewService() *Service {
	return &Service{
		validators: map[string]CountryValidator{
			"AT": &validators.Austria{},
			"BR": &validators.Brazil{},
			"GB": &validators.GB{},
			"BA": &validators.Bosnia{},
		},
	}
}

func (svc *Service) Parse(ibanStr string) (IBAN, error) {
	matches := ibanRegexp.FindStringSubmatch(ibanStr)
	if matches == nil || len(matches) != 4 {
		return IBAN{}, fmt.Errorf("provided string `%s` does not satisfy the iban format: %s", ibanStr, ibanRegexp.String())
	}

	countryCode, checkDigitsStr, bban := matches[1], matches[2], matches[3]
	checkDigits, err := strconv.ParseUint(checkDigitsStr, 10, 8)
	if err != nil {
		return IBAN{}, fmt.Errorf("failed to parse check-digits: %s for iban %s", ibanStr)
	}

	return IBAN{
		CountryCode: countryCode,
		CheckDigits: uint8(checkDigits),
		BBAN:        bban,
	}, nil
}

// Validate validates the iban's format and checks the check-digits.
func (svc *Service) Validate(i IBAN) error {
	if i.CountryCode == "" {
		return fmt.Errorf("country code is empty")
	}
	if i.BBAN == "" {
		return fmt.Errorf("bban is empty")
	}

	validator, ok := svc.validators[i.CountryCode]
	if !ok {
		return fmt.Errorf("country code %s is not supported", i.CountryCode)
	}

	err := validator.ValidateChecksum(i.String())
	if err != nil {
		return fmt.Errorf("check-digits are invalid: %s", err)
	}

	err = validator.ValidateLength(i.String())
	if err != nil {
		return fmt.Errorf("failed to validate length: %s", err)
	}

	err = validator.ValidateBban(i.BBAN)
	if err != nil {
		return fmt.Errorf("failed to validate bban: %s", err)
	}

	return nil
}
