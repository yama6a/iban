package iban

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_Parse(t *testing.T) {
	tests := []struct {
		name    string
		ibanStr string
		want    IBAN
		wantErr error
	}{
		{
			name:    "success",
			ibanStr: "NL02ABNA0123456789",
			want: IBAN{
				CountryCode: "NL",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
		},
		{
			name:    "fails for missing check digits",
			ibanStr: "NLXX0000001", // XX should be numeric
			wantErr: ErrIncorrectIbanFormat,
		},
		{
			name:    "fails for missing country code",
			ibanStr: "110000001", // should start with two uppercase letters
			wantErr: ErrIncorrectIbanFormat,
		},
		{
			name:    "fails for missing BBAN",
			ibanStr: "NL02", // should have digits after the check digits
			wantErr: ErrIncorrectIbanFormat,
		},
		{
			name:    "fails for empty string",
			ibanStr: "",
			wantErr: ErrIncorrectIbanFormat,
		},
	}
	for _, tt := range tests {
		svc := &Service{}

		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			got, err := svc.Parse(tt.ibanStr)
			require.Equal(t, tt.wantErr, err)
			if err != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name       string
		validators map[string]countryValidator
		i          IBAN
		wantErr    error
	}{
		{
			name: "success",
			i: IBAN{
				CountryCode: "NL",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
			validators: map[string]countryValidator{
				"NL": {
					CountryCode: "NL",
					Length:      18,
					BBANRegex:   regexp.MustCompile(`^[\dA-Z]+$`),
				},
			},
		},
		{
			name: "fails for missing country code",
			i: IBAN{
				CountryCode: "",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
			wantErr: ErrCountryCodeEmpty,
		},
		{
			name: "fails for missing BBAN",
			i: IBAN{
				CountryCode: "NL",
				CheckDigits: "02",
				BBAN:        "",
			},
			wantErr: ErrBBANEmpty,
		},
		{
			name: "fails for invalid country code",
			i: IBAN{
				CountryCode: "XX",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
			wantErr: ErrCountryCodeNotSupported,
		},
		{
			name: "fails for invalid length",
			i: IBAN{
				CountryCode: "NL",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
			validators: map[string]countryValidator{
				"NL": {
					CountryCode: "NL",
					Length:      5,
					BBANRegex:   regexp.MustCompile(`^[\dA-Z]+$`),
				},
			},
			wantErr: ErrIncorrectLength,
		},
		{
			name: "fails for invalid IBAN Checksum",
			i: IBAN{
				CountryCode: "GB",
				CheckDigits: "92", // should be 29
				BBAN:        "NWBK60161331926819",
			},
			validators: map[string]countryValidator{
				"GB": {
					CountryCode: "GB",
					Length:      22,
					BBANRegex:   regexp.MustCompile(`^[\dA-Z]+$`),
				},
			},
			wantErr: ErrIncorrectIBANChecksum,
		},
		{
			name: "fails invalid BBAN format",
			i: IBAN{
				CountryCode: "NL",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
			validators: map[string]countryValidator{
				"NL": {
					CountryCode: "NL",
					Length:      18,
					BBANRegex:   regexp.MustCompile(`^[\d]+$`), // only accepts digits
				},
			},
			wantErr: ErrIncorrectBBANFormat,
		},
		{
			name: "fails for incorrect BBAN checksum",
			i: IBAN{
				CountryCode: "NL",
				CheckDigits: "02",
				BBAN:        "ABNA0123456789",
			},
			validators: map[string]countryValidator{
				"NL": {
					CountryCode: "NL",
					Length:      18,
					BBANRegex:   regexp.MustCompile(`^[\dA-Z]+$`),
					BBANChecksumFunc: func(s string) bool {
						return false
					},
				},
			},
			wantErr: ErrIncorrectBBANChecksum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			svc := &Service{validators: tt.validators}
			require.ErrorIs(t, svc.Validate(tt.i), tt.wantErr)
		})
	}
}
