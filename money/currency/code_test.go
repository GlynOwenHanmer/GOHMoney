package currency_test

import (
	"testing"

	"encoding/json"

	"fmt"

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
		c, err := currency.NewCode(test.code)
		assert.Equal(t, test.err, err != nil)
		if err != nil {
			lenErr, ok := err.(currency.ErrInvalidCodeLength)
			assert.True(t, ok)
			assert.Equal(t, len(test.code), lenErr.Length)
		}
		assert.Equal(t, test.code, (*c).String())
	}
}

func TestJSON(t *testing.T) {
	ca, err := currency.NewCode("YEN")
	assert.Nil(t, err)
	bs, err := json.Marshal(ca)
	assert.Nil(t, err)
	cb, err := currency.UnmarshalJSON(bs)
	assert.Nil(t, err)
	assert.Equal(t, ca, cb)
}

func TestUnmarshalJSON(t *testing.T) {
	for _, test := range []struct {
		code string
		err  bool
	}{
		{code: "", err: true},
		{code: "YEN", err: false},
		{code: "QWERTYUIOP", err: true},
	} {
		json := fmt.Sprintf(`{"Code":"%s"}`, test.code)
		c, err := currency.UnmarshalJSON([]byte(json))
		assert.Equal(t, test.err, err != nil)
		if err != nil {
			lenErr, ok := err.(currency.ErrInvalidCodeLength)
			assert.True(t, ok)
			assert.Equal(t, len(test.code), lenErr.Length)
			continue
		}
		assert.Equal(t, test.code, (*c).String())
	}
}
