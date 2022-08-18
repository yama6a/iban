package iban

import (
	"fmt"
	"strconv"
)

type baseValidator struct{}

// mod97 calculates the mod97 of very large string-represented number (too large to fit into uint64)
//
//	We use Horner's Method for this (https://en.wikipedia.org/wiki/Horner%27s_method)
func mod97(s string) uint {
	var remainder uint
	for _, r := range s {
		remainder = (remainder * 10) + uint(r-'0')
		remainder %= 97
	}
	return remainder
}

// convertLetterToInt converts a letter to its corresponding number. (A = 10, B = 11, ..., Z = 35)
func convertLetterToInt(r rune) (int, error) {
	number, err := strconv.ParseInt(string(r), 36, 10)
	if err != nil {
		return 0, fmt.Errorf("failed to convert rune to Number: %s", err)
	}

	return int(number), nil
}
