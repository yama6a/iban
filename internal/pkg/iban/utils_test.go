package iban

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mod97(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want uint
	}{
		{
			name: "success three digits",
			s:    "123",
			want: 26,
		},
		{
			name: "success 10 digits",
			s:    "1234567890",
			want: 2,
		},
		{
			name: "success large number",
			s:    "3214282912345698765432161182",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			if got := mod97(tt.s); got != tt.want {
				t.Errorf("mod97() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertLetterToInt(t *testing.T) {
	tests := []struct {
		name    string
		r       rune
		want    int
		wantErr error
	}{
		{
			name: "success A",
			r:    'A',
			want: 10,
		},
		{
			name: "success B",
			r:    'B',
			want: 11,
		},
		{
			name: "success Z",
			r:    'Z',
			want: 35,
		},
		{
			name: "fails for invalid string",
			r:    '$',
			wantErr: fmt.Errorf(
				"failed to convert rune to Number: strconv.ParseInt: parsing \"$\": invalid syntax",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertLetterToInt(tt.r)
			require.Equal(t, err, tt.wantErr)
			if err != nil {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
