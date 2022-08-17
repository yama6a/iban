package iban

import "testing"

type mockParser struct {
	t            *testing.T
	ParseFunc    func(string) (IBAN, error)
	ValidateFunc func(IBAN) error
}

func (p *mockParser) Parse(s string) (IBAN, error) {
	if p.ParseFunc == nil {
		p.t.Fatalf("mockParser.ParseFunc: method is nil but Parser.Parse was just called")
	}
	return p.ParseFunc(s)
}

func (p *mockParser) Validate(i IBAN) error {
	if p.ValidateFunc == nil {
		p.t.Fatalf("mockParser.ValidateFunc: method is nil but Parser.Validate was just called")
	}
	return p.ValidateFunc(i)
}
