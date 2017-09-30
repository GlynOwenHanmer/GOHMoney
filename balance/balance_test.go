package balance_test

import (
	"errors"
	"testing"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney/balance"
)

func Test_ValidateBalance(t *testing.T) {
	invalidBalance := balance.Balance{}
	err := invalidBalance.Validate()
	if err != balance.BalanceZeroDate {
		t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", balance.BalanceZeroDate, err)
	}

	validBalance := balance.Balance{Date: time.Now()}
	err = validBalance.Validate()
	if err != nil {
		t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", error(nil), err)
	}
}

type BalanceErrorSet struct {
	balance.Balance
	error
}

func Test_Earliest_EmptyBalances(t *testing.T) {
	balances := balance.Balances{}
	expected := BalanceErrorSet{Balance: balance.Balance{}, error: errors.New(balance.EmptyBalancesMessage)}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithNoDate(t *testing.T) {
	balances := balance.Balances{balance.Balance{}}
	expected := BalanceErrorSet{balance.Balance{}, nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithSingleDate(t *testing.T) {
	earliest := balance.Balance{Date: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), Amount: 10}
	balances := balance.Balances{earliest}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest := balance.Balance{Date: date, Amount: 10}
	balances := balance.Balances{earliest, balance.Balance{Date: date, Amount: 20}}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest := balance.Balance{Date: date1, Amount: 10}
	balances := balance.Balances{balance.Balance{Date: date2}, earliest, balance.Balance{Date: date1, Amount: 20}, balance.Balance{Date: date3}}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func testEarliestSet(t *testing.T, expected BalanceErrorSet, balances balance.Balances) {
	actualBalance, actualError := balances.Earliest()
	actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
	testBalanceResults(t, expected, actual)
}

func Test_Latest_EmptyBalances(t *testing.T) {
	balances := balance.Balances{}
	expected := BalanceErrorSet{Balance: balance.Balance{}, error: errors.New(balance.EmptyBalancesMessage)}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithNoDate(t *testing.T) {
	balances := balance.Balances{balance.Balance{}}
	expected := BalanceErrorSet{balance.Balance{}, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSingleDate(t *testing.T) {
	latest := balance.Balance{Date: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), Amount: 10}
	balances := balance.Balances{latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	latest := balance.Balance{Date: date, Amount: 10}
	balances := balance.Balances{balance.Balance{Date: date, Amount: 20}, latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	latest := balance.Balance{Date: date3, Amount: 10}
	balances := balance.Balances{balance.Balance{Date: date2}, balance.Balance{Date: date3}, latest, balance.Balance{Date: date1, Amount: 20}}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func testLatestSet(t *testing.T, expected BalanceErrorSet, balances balance.Balances) {
	actualBalance, actualError := balances.Latest()
	actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
	testBalanceResults(t, expected, actual)
}

func testBalanceResults(t *testing.T, expected BalanceErrorSet, actual BalanceErrorSet) {
	if expected.error != actual.error {
		switch {
		case expected.error == nil:
			t.Errorf("Expected no error but got %v", actual)
		case actual.error == nil:
			t.Errorf("Error error (%v) but didn't get one", expected)
		case expected.error.Error() == actual.error.Error():
			break
		default:
			t.Errorf("Error unexpected\nExpected: %s\nActual  : %s", expected, actual)
		}
	}
	if expected.Balance != actual.Balance {
		t.Errorf("Balance unexpected\nExpected: %v\nActual  : %v", expected.Balance, actual.Balance)
	}
}

func TestBalances_Sum(t *testing.T) {
	testSets := []struct {
		balance.Balances
		expectedSum float32
	}{
		{
			expectedSum: 0,
		},
		{
			Balances: balance.Balances{
				balance.Balance{Amount: 1},
			},
			expectedSum: 1,
		},
		{
			Balances: balance.Balances{
				balance.Balance{Amount: 1},
				balance.Balance{Amount: 2},
			},
			expectedSum: 3,
		},

		{
			Balances: balance.Balances{
				balance.Balance{Amount: 1},
				balance.Balance{Amount: 2},
				balance.Balance{Amount: -3},
			},
			expectedSum: 0,
		},
	}

	for _, testSet := range testSets {
		actual := testSet.Balances.Sum()
		if testSet.expectedSum != actual {
			t.Errorf("Unexpected sum.\nExpected: %f\nActual  : %f\nBalances: %v", testSet.expectedSum, actual, testSet.Balances)
		}
	}
}
