package GOHMoney

import "testing"

func TestAccountFieldError_equal(t *testing.T) {
	testSets := []struct {
		errA, errB AccountFieldError
		equal      bool
	}{
		{
			errA:  AccountFieldError{},
			errB:  AccountFieldError{},
			equal: true,
		},
		{
			errA: AccountFieldError{
				EmptyNameError,
			},
			errB:  AccountFieldError{},
			equal: false,
		},
		{
			errA: AccountFieldError{
				EmptyNameError,
			},
			errB: AccountFieldError{
				EmptyNameError,
			},
			equal: true,
		},
		{
			errA: AccountFieldError{
				ZeroDateOpenedError,
				EmptyNameError,
			},
			errB: AccountFieldError{
				EmptyNameError,
				ZeroDateOpenedError,
			},
			equal: false,
		},
		{
			errA: AccountFieldError{
				ZeroDateOpenedError,
				EmptyNameError,
				ZeroValidDateClosedError,
			},
			errB: AccountFieldError{
				ZeroDateOpenedError,
				EmptyNameError,
				ZeroValidDateClosedError,
			},
			equal: true,
		},
	}
	for _, testSet := range testSets {
		equalA := testSet.errA.equal(testSet.errB)
		equalB := testSet.errB.equal(testSet.errA)
		if equalA != equalB {
			t.Fatalf("equal did not return same value when comparing AccountFieldError to other AccountFieldError the reverse way around.")
		}
		if testSet.equal != equalA {
			t.Errorf("Unexpected equal value.\n\tExpected: %t\n\tActual  : %t", testSet.equal, equalA)
		}
	}
}
