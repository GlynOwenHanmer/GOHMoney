package balance_test

import (
	"errors"
	"testing"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney/balance"
	"github.com/rhymond/go-money"
)

func TestNew(t *testing.T) {
	var expected error

	expected = balance.ZeroDate
	invalidTime := time.Time{}
	if _, err := balance.New(invalidTime, balance.NewMoney(0)); err != expected {
		t.Errorf("Unexpected error\nExpected: %s\nActual  : %s", expected, err)
	}

	expected = nil
	validTime := time.Now()
	if _, err := balance.New(validTime, balance.NewMoney(0)); err != expected {
		t.Errorf("Unexpected error\nExpected: %s\nActual  : %s", expected, err)
	}

	expected = nil
	if _, err := balance.New(validTime, money.Money{}); err != expected {
		t.Errorf("Unexpected error\nExpected: %s\nActual  : %s", expected, err)
	}
}

func Test_ValidateBalance(t *testing.T) {
	invalidBalance := balance.Balance{}
	err := invalidBalance.Validate()
	if err != balance.ZeroDate {
		t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", balance.ZeroDate, err)
	}

	validBalance, err := balance.New(time.Now(), balance.NewMoney(0))
	if err != nil {
		t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", error(nil), err)
	}
	if err := validBalance.Validate(); err != nil {
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
	earliest, _ := balance.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.NewMoney(10))
	balances := balance.Balances{earliest}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest, _ := balance.New(date, balance.NewMoney(10))
	other, _ := balance.New(date, balance.NewMoney(20))
	balances := balance.Balances{earliest, other}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest, _ := balance.New(date1, balance.NewMoney(10))
	other, _ := balance.New(date2, balance.NewMoney(0))
	other2, _ := balance.New(date1, balance.NewMoney(20))
	other3, _ := balance.New(date3, balance.NewMoney(489))
	balances := balance.Balances{other, earliest, other2, other3}
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
	latest, _ := balance.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), balance.NewMoney(10))
	balances := balance.Balances{latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	latest, _ := balance.New(date, balance.NewMoney(10))
	other, _ := balance.New(date, balance.NewMoney(20))
	balances := balance.Balances{other, latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	latest, _ := balance.New(date3, balance.NewMoney(20))
	other1, _ := balance.New(date2, balance.NewMoney(0))
	other2, _ := balance.New(date3, balance.NewMoney(0))
	other3, _ := balance.New(date1, balance.NewMoney(20))
	balances := balance.Balances{other1, other2, latest, other3}
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
		amounts     []int64
		expectedSum money.Money
	}{
		{
			expectedSum: balance.NewMoney(0),
		},
		{
			amounts:     []int64{1},
			expectedSum: balance.NewMoney(1),
		},
		{
			amounts:     []int64{1, 2},
			expectedSum: balance.NewMoney(3),
		},

		{
			amounts:     []int64{1, 2, -3},
			expectedSum: balance.NewMoney(0),
		},
	}

	now := time.Now()

	for _, testSet := range testSets {
		var bs balance.Balances
		for _, a := range testSet.amounts {
			b, _ := balance.New(now, balance.NewMoney(a))
			bs = append(bs, b)
		}
		actual, err := bs.Sum()
		if err != nil {
			t.Fatalf("Error summing balances: %s", err)
		}
		equal, err := actual.Equals(&testSet.expectedSum)
		if !equal {
			t.Errorf("Unexpected sum.\nExpected: %f\nActual  : %f\nBalances: %v", testSet.expectedSum, actual, bs)
		}
	}
}