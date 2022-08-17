package iban

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestController_writeResponse(t *testing.T) {
	type args struct {
		iban   *IBAN
		err    error
		status int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				iban:   &IBAN{CountryCode: "NL", CheckDigits: "12", BBAN: "112233"},
				err:    nil,
				status: http.StatusOK,
			},
			want: `{"error":null,"is_valid":true,"iban":{"country_code":"NL","check_digits":"12","bban":"112233"}}`,
		},
		{
			name: "error with IBAN object",
			args: args{
				iban:   &IBAN{CountryCode: "NL", CheckDigits: "12", BBAN: "112233"},
				err:    errors.New("some error"),
				status: http.StatusTeapot,
			},
			want: `{"error":"some error","is_valid":false,"iban":{"country_code":"NL","check_digits":"12","bban":"112233"}}`,
		},
		{
			name: "error without IBAN object",
			args: args{
				iban:   nil,
				err:    errors.New("other error"),
				status: http.StatusPreconditionFailed,
			},
			want: `{"error":"other error","is_valid":false,"iban":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // shadow tt for parallel execution
			t.Parallel()

			crtl := Controller{parser: nil, logger: zap.NewNop()}
			w := httptest.NewRecorder()
			crtl.writeResponse(w, tt.args.iban, tt.args.err, tt.args.status)
			require.JSONEq(t, tt.want, w.Body.String())
		})
	}
}

func TestController_validate(t *testing.T) {
	tests := []struct {
		name               string
		r                  *http.Request
		parser             Parser
		want               string
		expectedStatusCode int
	}{
		{
			name: "success",
			r:    httptest.NewRequest(http.MethodGet, "/v1/iban/NL22555566667777/validate", nil),
			want: `{"error":null,"is_valid":true,"iban":{"country_code":"NL","check_digits":"22","bban":"555566667777"}}`,
			parser: &mockParser{
				ParseFunc: func(s string) (IBAN, error) {
					return IBAN{CountryCode: "NL", CheckDigits: "22", BBAN: "555566667777"}, nil
				},
				ValidateFunc: func(iban IBAN) error {
					return nil
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "parsing error returns 422",
			r:    httptest.NewRequest(http.MethodGet, "/v1/iban/NL22555566667777/validate", nil),
			want: fmt.Sprintf(`{"error":"%s","is_valid":false,"iban":null}`, "parsing error"),
			parser: &mockParser{
				ParseFunc: func(s string) (IBAN, error) {
					return IBAN{}, errors.New("parsing error")
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "validation error returns 422",
			r:    httptest.NewRequest(http.MethodGet, "/v1/iban/NL22555566667777/validate", nil),controller
			want: fmt.Sprintf(`{"error":"%s","is_valid":false,"iban":{"country_code":"NL","check_digits":"22","bban":"555566667777"}}`, "validation error"),
			parser: &mockParser{
				ParseFunc: func(s string) (IBAN, error) {
					return IBAN{CountryCode: "NL", CheckDigits: "22", BBAN: "555566667777"}, nil
				},
				ValidateFunc: func(iban IBAN) error {
					return errors.New("validation error")
				},
			},
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // shadow tt for parallel execution
			t.Parallel()

			ctrl := Controller{parser: tt.parser, logger: zap.NewNop()}

			w := httptest.NewRecorder()
			ctrl.validate(w, tt.r)
			require.JSONEq(t, tt.want, w.Body.String())
		})
	}
}
