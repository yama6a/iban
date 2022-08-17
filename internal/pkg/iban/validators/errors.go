package validators

import "errors"

var (
	ErrIncorrectLength       = errors.New("IBAN has the incorrect length for the specified country")
	ErrIncorrectBBANFormat   = errors.New("IBAN has the incorrect BBAN format for the specified country")
	ErrIncorrectBBANChecksum = errors.New("IBAN has the incorrect BBAN checksum for the specified country")
	ErrIncorrectIBANChecksum = errors.New("IBAN has the incorrect checksum")
)
