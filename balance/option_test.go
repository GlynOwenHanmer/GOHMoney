package balance_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glynternet/GOHMoney/balance"
	"github.com/glynternet/GOHMoney/common"
	"github.com/stretchr/testify/assert"
)

func TestAmount(t *testing.T) {
	b, err := balance.New(time.Now())
	common.FatalIfError(t, err, "Creating Balance")
	assert.Equal(t, int64(0), b.Amount())
	assert.Nil(t, balance.Amount(-645)(&b))
	assert.Equal(t, int64(-645), b.Amount())
}

func TestErrorOption(t *testing.T) {
	errorFn := func(a *balance.Balance) error {
		return errors.New("TEST ERROR")
	}
	_, err := balance.New(time.Now(), errorFn)
	assert.Equal(t, errors.New("TEST ERROR"), err)
}

func TestCurrencyCode(t *testing.T) {

}
