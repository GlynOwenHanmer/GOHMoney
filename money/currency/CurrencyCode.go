package currency

import "fmt"

const InvalidCodeIdentifier = "INVALID CODE"

func New(c string) (Code, error) {
	code := Code(c)
	err := code.Validate()
	if err != nil {
		return Code(InvalidCodeIdentifier), err
	}
	return code, nil
}

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

func (c Code)Validate() error {
	err := validateCodeLength(c)
	if err != nil {
		return err
	}
	return nil
}