package account_test

import (
	"testing"
	"time"

	"github.com/glynternet/GOHMoney/account"
	"github.com/glynternet/GOHMoney/common"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestClosedTime(t *testing.T) {
	start := time.Now()
	closeA := start.Add(72 * time.Hour)
	closeFn := account.CloseTime(closeA)
	a, err := account.New("TEST_ACCOUNT", start, closeFn)
	common.FatalIfError(t, err, "Creating Account")
	assert.True(t, a.End().EqualTime(closeA))

	closeB := closeA.Add(100 * time.Hour)
	err = account.CloseTime(closeB)(&a)
	common.FatalIfError(t, err, "Creating CloseTime Option")
	assert.True(t, a.End().EqualTime(closeB))
}

func TestErrorOption(t *testing.T) {
	errorFn := func(a *account.Account) error {
		return errors.New("TEST ERROR")
	}
	_, err := account.New("TEST_ACCOUNT", time.Now(), errorFn)
	assert.Equal(t, errors.New("TEST ERROR"), err)
}