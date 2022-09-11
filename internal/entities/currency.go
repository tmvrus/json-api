package entities

import "fmt"

type Currency string

func (c Currency) Valid() error {
	const ISO4217Length = 3

	if len(c) != ISO4217Length {
		return fmt.Errorf("currency should be three-character ISO 4217 code")
	}
	return nil
}
