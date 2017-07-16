package GOHMoney

import (
	"testing"
	"github.com/lib/pq"
	"time"
)

func Test_ValidateAccount(t *testing.T) {
	testSets := []struct{
		insertedAccount Account
		AccountFieldError
	}{
		{
			insertedAccount:   Account{},
			AccountFieldError: AccountFieldError{EmptyNameError, ZeroDateOpenedError},
		},
		{
			insertedAccount: Account{
				Name:"TEST_ACCOUNT",
			},
			AccountFieldError: AccountFieldError{ZeroDateOpenedError},
		},
		{
			insertedAccount: Account{
				TimeRange:TimeRange{
					Start:pq.NullTime{
						Valid: true,
						Time: time.Date(2000,1,1,1,1,1,1,time.UTC),
					},
				},
			},
			AccountFieldError: AccountFieldError{EmptyNameError},
		},
		{
			insertedAccount: Account{
				Name:"TEST_ACCOUNT",
				TimeRange: TimeRange{
					Start:pq.NullTime{},
					End:pq.NullTime{
						Valid:true,
					},
				},
			},
			AccountFieldError: AccountFieldError{ZeroDateOpenedError, ZeroValidDateClosedError},
		},
		{
			insertedAccount: Account{
				Name:"TEST_ACCOUNT",
				TimeRange: TimeRange{
					Start:pq.NullTime{
						Valid:true,
						Time:time.Date(2000,1,1,1,1,1,1,time.UTC),
					},
					End:pq.NullTime{
						Valid:true,
						Time:time.Date(1999,1,1,1,1,1,1,time.UTC),
					},
				},
			},
			AccountFieldError: AccountFieldError{string(DateClosedBeforeDateOpenedError)},
		},
		{
			insertedAccount:   newTestAccount(),
			AccountFieldError: nil,
		},
	}
	for _, testSet := range testSets {
		actual := testSet.insertedAccount.Validate()
		expected := testSet.AccountFieldError
		if !stringSlicesMatch(expected, actual) {
			t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s\nInserted Account: %s", expected, actual, testSet.insertedAccount)
		}
	}
}

func newTestAccount() Account {
	return Account{
		Name:       "TEST_ACCOUNT",
		TimeRange: TimeRange{
			Start:pq.NullTime{
				Valid: true,
				Time: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			},
			End:pq.NullTime{
				Valid: true,
				Time: time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC),
			},
		},
	}
}

func stringSlicesMatch(array1, array2 []string) bool {
	if len(array1) != len(array2) {
		return false
	}
	for i := 0; i < len(array1); i++ {
		if array1[i] != array2[i] {
			return false
		}
	}
	return true
}

func Test_IsOpen(t *testing.T) {
	testSets := []struct {
		Account
		IsOpen bool
	}{
		{
			Account: Account{},
			IsOpen:  true,
		},
		{
			Account: Account{
				TimeRange:TimeRange{
					End:pq.NullTime{Valid:false},
				},
			},
			IsOpen:  true,
		},
		{
			Account: Account{
				TimeRange: TimeRange{
					End: pq.NullTime{Valid: true},
				},
			},
			IsOpen:  false,
		},
	}
	for _, testSet := range testSets {
		actual := testSet.Account.IsOpen()
		if actual != testSet.IsOpen {
			t.Errorf("Account IsOpen expected %t, got %t. Account: %v", testSet.IsOpen, actual, testSet.Account)
		}
	}
}