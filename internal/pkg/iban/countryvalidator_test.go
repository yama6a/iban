package iban

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_countryValidator_ValidateIbanLength(t *testing.T) {
	tests := []struct {
		name    string
		Length  int
		iban    IBAN
		wantErr error
	}{
		{
			name:   "valid length of 10",
			Length: 10,
			iban:   IBAN{CountryCode: "BR", CheckDigits: "12", BBAN: "123456"},
		},
		{
			name:   "invalid length of 14",
			Length: 14,
			iban:   IBAN{CountryCode: "BR", CheckDigits: "12", BBAN: "1234567890"},
		},
		{
			name:   "invalid length of 20 with letters",
			Length: 20,
			iban:   IBAN{CountryCode: "BR", CheckDigits: "12", BBAN: "1234567890ABCDEF"},
		},
		{
			name:    "too short returns error",
			Length:  10,
			iban:    IBAN{CountryCode: "BR", CheckDigits: "12", BBAN: "1234"}, // only 8 instead of 10
			wantErr: ErrIncorrectLength,
		},
		{
			name:    "too long returns error",
			Length:  10,
			iban:    IBAN{CountryCode: "BR", CheckDigits: "12", BBAN: "12345678"}, // 12 instead of 10
			wantErr: ErrIncorrectLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			c := countryValidator{Length: tt.Length}
			err := c.ValidateIbanLength(tt.iban)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_countryValidator_ValidateBbanFormat(t *testing.T) {
	tests := []struct {
		name      string
		BBANRegex *regexp.Regexp
		iban      IBAN
		wantErr   error
	}{
		{
			name:      "valid BBAN",
			BBANRegex: regexp.MustCompile(`^\d{10}[A-Z0-9]{11}\d{2}$`),
			iban:      IBAN{CountryCode: "FR", CheckDigits: "12", BBAN: "1234567890ABC12345DEF01"},
		},
		{
			name:      "invalid BBAN",
			BBANRegex: regexp.MustCompile(`^\d{10}[A-Z0-9]{11}\d{2}$`),
			iban:      IBAN{CountryCode: "FR", CheckDigits: "12", BBAN: "1234567890ABC12345DEF0"}, // missing last digit
			wantErr:   ErrIncorrectBBANFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			c := countryValidator{BBANRegex: tt.BBANRegex}
			err := c.ValidateBbanFormat(tt.iban)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_countryValidator_ValidateBbanChecksum(t *testing.T) {
	tests := []struct {
		name             string
		BBANChecksumFunc func(string) bool
		iban             IBAN
		wantErr          error
	}{
		{
			name:             "valid BBAN",
			BBANChecksumFunc: func(bban string) bool { return true },
			iban:             IBAN{CountryCode: "FR", CheckDigits: "12", BBAN: "1234567890ABC12345DEF01"},
		},
		{
			name:             "invalid BBAN",
			BBANChecksumFunc: func(bban string) bool { return false },
			iban:             IBAN{CountryCode: "FR", CheckDigits: "12", BBAN: "1234567890ABC12345DEF01"},
			wantErr:          ErrIncorrectBBANChecksum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			c := countryValidator{BBANChecksumFunc: tt.BBANChecksumFunc}
			err := c.ValidateBbanChecksum(tt.iban)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_countryValidator_ValidateIbanChecksum(t *testing.T) {
	tests := []struct {
		name    string
		iban    IBAN
		wantErr error
	}{
		{
			name: "valid IBAN",
			iban: IBAN{CountryCode: "GB", CheckDigits: "29", BBAN: "NWBK60161331926819"},
		},
		{
			name:    "invalid IBAN - wrong check digits",
			iban:    IBAN{CountryCode: "GB", CheckDigits: "92", BBAN: "NWBK60161331926819"},
			wantErr: ErrIncorrectIBANChecksum,
		},
		{
			name:    "invalid IBAN - wrong country code",
			iban:    IBAN{CountryCode: "US", CheckDigits: "29", BBAN: "NWBK60161331926819"},
			wantErr: ErrIncorrectIBANChecksum,
		},
		{
			name:    "invalid IBAN - wrong BBAN",
			iban:    IBAN{CountryCode: "GB", CheckDigits: "29", BBAN: "NWBK60161331926810"}, // last digit should be 9
			wantErr: ErrIncorrectIBANChecksum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			c := countryValidator{}
			err := c.ValidateIbanChecksum(tt.iban)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
