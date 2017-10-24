package account

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/glynternet/GOHMoney/balance"
	"github.com/glynternet/GOHMoney/common"
	"github.com/glynternet/GOHMoney/money/currency"
	gtime "github.com/glynternet/go-time"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateAccount(t *testing.T) {
	testSets := []struct {
		insertedAccount Account
		FieldError
	}{
		{
			insertedAccount: Account{},
			FieldError:      FieldError{EmptyNameError, ZeroDateOpenedError},
		},
		{
			insertedAccount: Account{
				Name: "TEST_ACCOUNT",
			},
			FieldError: FieldError{ZeroDateOpenedError},
		},
		{
			insertedAccount: Account{
				timeRange: gtime.Range{
					Start: gtime.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
					},
				},
			},
			FieldError: FieldError{EmptyNameError},
		},
		{
			insertedAccount: Account{
				Name: "TEST_ACCOUNT",
				timeRange: gtime.Range{
					Start: gtime.NullTime{},
					End: gtime.NullTime{
						Valid: true,
					},
				},
			},
			FieldError: FieldError{ZeroDateOpenedError, ZeroValidDateClosedError},
		},
		{
			insertedAccount: Account{
				Name: "TEST_ACCOUNT",
				timeRange: gtime.Range{
					Start: gtime.NullTime{
						Valid: true,
						Time:  time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
					},
					End: gtime.NullTime{
						Valid: true,
						Time:  time.Date(1999, 1, 1, 1, 1, 1, 1, time.UTC),
					},
				},
			},
			FieldError: FieldError{string(gtime.EndTimeBeforeStartTime)},
		},
		{
			insertedAccount: newTestAccount(),
			FieldError:      nil,
		},
		{
			insertedAccount: Account{
				Name: "TEST_ACCOUNT",
				timeRange: gtime.Range{
					Start: gtime.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
					End: gtime.NullTime{},
				},
			},
			FieldError: nil,
		},
	}
	for _, testSet := range testSets {
		actual := testSet.insertedAccount.Validate()
		expected := testSet.FieldError
		if !stringSlicesMatch(expected, actual) {
			t.Errorf("Unexpected error.\nExpected: %+v\nActual  : %+v\nInserted Account: %+v", expected, actual, testSet.insertedAccount)
		}
	}
}

func newTestAccount() Account {
	return Account{
		Name: "TEST_ACCOUNT",
		timeRange: gtime.Range{
			Start: gtime.NullTime{
				Valid: true,
				Time:  time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			},
			End: gtime.NullTime{
				Valid: true,
				Time:  time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC),
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
				timeRange: gtime.Range{
					End: gtime.NullTime{Valid: false},
				},
			},
			IsOpen: true,
		},
		{
			Account: Account{
				timeRange: gtime.Range{
					End: gtime.NullTime{Valid: true},
				},
			},
			IsOpen: false,
		},
	}
	for _, testSet := range testSets {
		actual := testSet.Account.IsOpen()
		if actual != testSet.IsOpen {
			t.Errorf("Account IsOpen expected %t, got %t. Account: %v", testSet.IsOpen, actual, testSet.Account)
		}
	}
}

func Test_InvalidAccountValidateBalance(t *testing.T) {
	present := time.Now()
	past := present.AddDate(-1, 0, 0)
	future := present.AddDate(1, 0, 0)

	invalidAccount := Account{
		timeRange: gtime.Range{
			Start: gtime.NullTime{
				Valid: true,
				Time:  future,
			},
			End: gtime.NullTime{
				Valid: true,
				Time:  past,
			},
		},
	}

	accountErr := invalidAccount.Validate()
	balanceErr := invalidAccount.ValidateBalance(balance.Balance{})
	switch {
	case accountErr == nil && balanceErr == nil:
		break
	case accountErr == nil || balanceErr == nil,
		accountErr.Error() != balanceErr.Error():
		t.Errorf("ValidateBalance did not return Account error when given invalid account.\nExpected: %s\nActual  : %s", accountErr, balanceErr)
	}
}

func Test_AccountValidateBalance(t *testing.T) {
	present := time.Now()
	past := present.AddDate(-1, 0, 0)
	future := present.AddDate(1, 0, 0)

	openAccount := Account{
		Name: "Test Account",
		timeRange: gtime.Range{
			Start: gtime.NullTime{
				Valid: true,
				Time:  present,
			},
			End: gtime.NullTime{Valid: false},
		},
	}
	closedAccount := Account{
		Name: "Test Account",
		timeRange: gtime.Range{
			Start: gtime.NullTime{
				Valid: true,
				Time:  present,
			},
			End: gtime.NullTime{
				Valid: true,
				Time:  future,
			},
		},
	}

	pastBalance, _ := balance.New(past, balance.Amount(1))
	presentBalance, _ := balance.New(present, balance.Amount(98737879))
	futureBalance, _ := balance.New(future, balance.Amount(-9876))
	evenFuturerBalance, _ := balance.New(future.AddDate(1, 0, 0), balance.Amount(-987654))
	testSets := []struct {
		Account
		balance.Balance
		error
	}{
		{
			Account: openAccount,
		},
		{
			Account: openAccount,
			Balance: pastBalance,
			error: balance.DateOutOfAccountTimeRange{
				BalanceDate:      pastBalance.Date(),
				AccountTimeRange: openAccount.timeRange,
			},
		},
		{
			Account: openAccount,
			Balance: presentBalance,
			error:   nil,
		},
		{
			Account: openAccount,
			Balance: futureBalance,
			error:   nil,
		},
		{
			Account: closedAccount,
			Balance: pastBalance,
			error: balance.DateOutOfAccountTimeRange{
				BalanceDate:      pastBalance.Date(),
				AccountTimeRange: closedAccount.timeRange,
			},
		},
		{
			Account: closedAccount,
			Balance: presentBalance,
			error:   nil,
		},
		{
			Account: closedAccount,
			Balance: futureBalance,
			error:   nil,
		},
		{
			Account: closedAccount,
			Balance: evenFuturerBalance,
			error: balance.DateOutOfAccountTimeRange{
				BalanceDate:      futureBalance.Date().AddDate(1, 0, 0),
				AccountTimeRange: closedAccount.timeRange,
			},
		},
	}
	for _, testSet := range testSets {
		err := testSet.Account.ValidateBalance(testSet.Balance)
		if testSet.error == err {
			continue
		}
		testSetTyped, testSetErrorIsType := testSet.error.(balance.DateOutOfAccountTimeRange)
		actualErrorTyped, actualErrorIsType := err.(balance.DateOutOfAccountTimeRange)
		if testSetErrorIsType != actualErrorIsType {
			t.Errorf("Expected and resultant errors are differently typed.\nExpected: %s\nActual  : %s", testSet.error, err)
			t.Logf("Account: %v\nbalance: %v", testSet.Account, testSet.Balance)
			continue
		}
		switch {
		case testSetTyped.AccountTimeRange.Equal(actualErrorTyped.AccountTimeRange):
			fallthrough
		case testSetTyped.BalanceDate.Equal(actualErrorTyped.BalanceDate):
			continue
		}
		var message bytes.Buffer
		fmt.Fprintf(&message, "Unexpected error.\nExpected: %+v\nActual  : %+v", testSetTyped, actualErrorTyped)
		fmt.Fprintf(&message, "\nExpected error: BalanceDate: %s, Accountgohtime.Range: %+v", testSetTyped.BalanceDate, testSetTyped.AccountTimeRange)
		fmt.Fprintf(&message, "\nActual error  : BalanceDate: %s, Accountgohtime.Range: %+v", actualErrorTyped.BalanceDate, actualErrorTyped.AccountTimeRange)
		t.Errorf(message.String())
	}
}

func Test_NewAccount(t *testing.T) {
	now := time.Now()
	type testSet struct {
		name  string
		start time.Time
		end   time.Time
		error
	}
	testSets := []testSet{
		{
			name:  "TEST_ACCOUNT",
			start: now,
		},
		{
			name:  "TEST_ACCOUNT_WITH_ACCOUNT_ERROR",
			start: now,
			end:   now.AddDate(0, 0, -1),
			error: FieldError{gtime.EndTimeBeforeStartTime.Error()},
		},
		{
			name:  "TEST_ACCOUNT",
			start: now,
			end:   now.AddDate(0, 0, +1),
		},
	}
	logTestSet := func(ts testSet) { t.Logf("Start: %s,\tEnd: %v,", ts.start, ts.end) }
	for _, set := range testSets {
		close := CloseTime(set.end)
		a, err := New(set.name, newTestCurrency(t, "YEN"), set.start, close)
		if !testNewAccountErrorTypes(t, set.error, err) {
			logTestSet(set)
		}
		if a.Name != set.name {
			t.Errorf("Unexpected name.\n\tExpected: %s\n\tActual  : %s", set.name, a.Name)
			logTestSet(set)
		}
		if !a.timeRange.Start.EqualTime(set.start) {
			t.Errorf("Unexpected start.\n\tExpected: %s\n\tActual  : %s", set.start, a.Start())
			logTestSet(set)
		}
		switch {
		case !set.end.IsZero():
			assert.True(t, a.End().EqualTime(set.end), "End NullTime should be equal to set.end when set.end is non Zero")
		case set.end.IsZero():
			assert.False(t, a.End().Valid, "End should not be Valid when set.end IsZero")
		}
	}
}

func testNewAccountErrorTypes(t *testing.T, expected, actual error) bool {
	expectedFieldError, expectedIsTyped := expected.(FieldError)
	actualFieldError, actualIsTyped := actual.(FieldError)
	switch {
	case actualIsTyped != expectedIsTyped,
		!actualIsTyped && actual != expected:
		t.Errorf("Unexpected error.\n\tExpected: %s\n\tActual  : %s", expected, actual)
		return false
	case actualIsTyped && !actualFieldError.Equal(expectedFieldError):
		t.Errorf("Error is correct type but unexpected contents.\n\tExpected: %s\n\tActual  : %s", expectedFieldError, actualFieldError)
		return false
	}
	return true
}

func newTestCurrency(t *testing.T, code string) currency.Code {
	c, err := currency.New(code)
	common.FatalIfError(t, err, "Creating New Currency Code")
	return c
}
