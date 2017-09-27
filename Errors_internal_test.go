package GOHMoney_test

import (
	"testing"
	"github.com/GlynOwenHanmer/GOHMoney"
)

func TestAccountFieldError_equal(t *testing.T) {
	testSets := []struct {
		errA, errB GOHMoney.AccountFieldError
		equal      bool
	}{
		{
			errA:  GOHMoney.AccountFieldError{},
			errB:  GOHMoney.AccountFieldError{},
			equal: true,
		},
		{
			errA: GOHMoney.AccountFieldError{
				GOHMoney.EmptyNameError,
			},
			errB:  GOHMoney.AccountFieldError{},
			equal: false,
		},
		{
			errA: GOHMoney.AccountFieldError{
				GOHMoney.EmptyNameError,
			},
			errB: GOHMoney.AccountFieldError{
				GOHMoney.EmptyNameError,
			},
			equal: true,
		},
		{
			errA: GOHMoney.AccountFieldError{
				GOHMoney.ZeroDateOpenedError,
				GOHMoney.EmptyNameError,
			},
			errB: GOHMoney.AccountFieldError{
				GOHMoney.EmptyNameError,
				GOHMoney.ZeroDateOpenedError,
			},
			equal: false,
		},
		{
			errA: GOHMoney.AccountFieldError{
				GOHMoney.ZeroDateOpenedError,
				GOHMoney.EmptyNameError,
				GOHMoney.ZeroValidDateClosedError,
			},
			errB: GOHMoney.AccountFieldError{
				GOHMoney.ZeroDateOpenedError,
				GOHMoney.EmptyNameError,
				GOHMoney.ZeroValidDateClosedError,
			},
			equal: true,
		},
	}
	for _, testSet := range testSets {
		equalA := testSet.errA.Equal(testSet.errB)
		equalB := testSet.errB.Equal(testSet.errA)
		if equalA != equalB {
			t.Fatalf("Equal did not return same value when comparing GOHMoney.AccountFieldError to other GOHMoney.AccountFieldError the reverse way around.")
		}
		if testSet.equal != equalA {
			t.Errorf("Unexpected Equal value.\n\tExpected: %t\n\tActual  : %t", testSet.equal, equalA)
		}
	}
}
