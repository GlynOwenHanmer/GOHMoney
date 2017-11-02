package currency

import "fmt"

// InvalidCodeIdentifier is the containing string of a Code that is created
// when an invalid string is given to New
const InvalidCodeIdentifier = "INVALID CODE"

// New returns a new Code if a valid string is given.
func New(c string) (code Code, err error) {
	code = Code(c)
	err = code.Validate()
	if err != nil {
		code = Code(InvalidCodeIdentifier)
	}
	return
}

// Code is a 3 character string representing a code for a currency
type Code string

func invalidCodeLength(length int) error {
	return fmt.Errorf("invalid currency code length (%d)", length)
}

func validateCodeLengthError(code string) (err error) {
	if length := len(code); length != 3 {
		err = invalidCodeLength(length)
	}
	return
}

// Validate returns an error if a Code is invalid
func (c Code) Validate() error {
	return validateCodeLengthError(string(c))
}
