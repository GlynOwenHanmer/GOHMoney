package currency

import "fmt"

// InvalidCodeIdentifier is the cotaining string of a Code that is created
// when an invalid string is given to New
const InvalidCodeIdentifier = "INVALID CODE"

// New returns a new Code if a valid string is given.
func New(c string) (Code, error) {
	code := Code(c)
	err := code.Validate()
	if err != nil {
		return Code(InvalidCodeIdentifier), err
	}
	return code, nil
}

// Code is a 3 character string representing a code for a currency
type Code string

func invalidCodeLength(length int) error {
	return fmt.Errorf("Invalid currency code length (%d)", length)
}

func validateCodeLength(code string) error {
	if length := len(code); length != 3 {
		return invalidCodeLength(length)
	}
	return nil
}

// Validate returns an error if a Code is invalid
func (c Code)Validate() error {
	return validateCodeLength(string(c))
}