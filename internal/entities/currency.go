package entities

import "fmt"

type Currency string

func (c Currency) Valid() error {
	if len(c) != 3 {
		return fmt.Errorf("currency should be three-character ISO 4217 code")
	}
	return nil
}
