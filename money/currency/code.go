package currency

import (
	"encoding/json"
	"fmt"
)

// NewCode returns a new code if a valid string is given.
func NewCode(currencyCode string) (c *Code, err error) {
	c = new(Code)
	*c = code(currencyCode)
	err = (*c).(code).validate()
	if err != nil {
		c = nil
	}
	return
}

type Code interface {
	String() string
}

// code is a 3 character string representing a code for a currency
type code string

func (c code) String() string {
	return string(c)
}

func (c code) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func UnmarshalJSON(data []byte) (*Code, error) {
	var aux string
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return nil, err
	}
	return NewCode(aux)
}

// validate returns an error if a code is invalid
func (c code) validate() error {
	return validateCodeLengthError(string(c))
}

func validateCodeLengthError(code string) (err error) {
	if length := len(code); length != 3 {
		err = ErrInvalidCodeLength{length}
	}
	return
}

type ErrInvalidCodeLength struct {
	Length int
}

func (e ErrInvalidCodeLength) Error() string {
	return fmt.Sprintf("invalid currency code Length (%d)", e.Length)
}
