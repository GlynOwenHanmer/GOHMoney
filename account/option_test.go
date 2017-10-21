package account_test

import (
	"testing"
	"github.com/glynternet/GOHMoney/account"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/glynternet/GOHMoney/common"
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