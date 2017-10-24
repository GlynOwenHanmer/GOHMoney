package currency_test

import (
	"testing"

	"github.com/glynternet/GOHMoney/money/currency"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		code string
		err  bool
	}{
		{code: "", err: true},
		{code: "YEN", err: false},
		{code: "QWERTYUIOP", err: true},
	} {
		c, err := currency.New(test.code)
		assert.Equal(t, test.err, err != nil)
		if err != nil {
			test.code = currency.InvalidCodeIdentifier
		}
		assert.Equal(t, test.code, string(c))
	}
}
