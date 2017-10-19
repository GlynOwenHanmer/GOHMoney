package balance_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/glynternet/GOHMoney/balance"
	"github.com/glynternet/GOHMoney/common"
	"github.com/glynternet/GOHMoney/money"
	innermoney "github.com/rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var expected error = balance.ZeroDate
	invalidTime := time.Time{}
	m, _ := money.New(0, "EUR")
	if _, err := balance.New(invalidTime, *m); err != expected {
		t.Errorf("Unexpected error\nExpected: %s\nActual  : %s", expected, err)
	}

	expected = nil
	validTime := time.Now()
	m, _ = money.New(0, "GBP")
	if _, err := balance.New(validTime, *m); err != expected {
		t.Errorf("Unexpected error\nExpected: %s\nActual  : %s", expected, err)
	}

	expected = money.ErrNoCurrency
	if _, err := balance.New(validTime, money.Money{}); err != expected {
		t.Errorf("Unexpected error\nExpected: %s\nActual  : %s", expected, err)
	}
}

func TestValidateBalance(t *testing.T) {
	invalidBalance := balance.Balance{}
	err := invalidBalance.Validate()
	if err != balance.ZeroDate {
		t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", balance.ZeroDate, err)
	}
	money, err := money.New(12, "GBP")
	common.FatalIfError(t, err, "Creating Money for testing")
	validBalance, err := balance.New(time.Now(), *money)
	common.ErrorIfError(t, err, "Creating Balance for testing")
	if err := validBalance.Validate(); err != nil {
		t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", error(nil), err)
	}
}

func TestBalance_Equal(t *testing.T) {
	newBalance := func(t *testing.T, date time.Time, amount int64, currency string) balance.Balance {
		m, err := money.New(amount, currency)
		common.FatalIfError(t, err, "Creating new Money")
		b, err := balance.New(date, *m)
		common.FatalIfError(t, err, "Creating new Balance")
		return b
	}

	now := time.Now()
	a := newBalance(t, now, 123, "GBP")
	b := newBalance(t, now, 123, "GBP")
	assert.True(t, a.Equal(b))

	b = newBalance(t, now, -123, "GBP")
	assert.True(t, !a.Equal(b))

	b = newBalance(t, now, 123, "EUR")
	assert.True(t, !a.Equal(b))

	b = newBalance(t, now.Add(1), 123, "GBP")
	assert.True(t, !a.Equal(b))
}

type BalanceErrorSet struct {
	balance.Balance
	error
}

func TestBalances_Earliest_EmptyBalances(t *testing.T) {
	balances := balance.Balances{}
	expected := BalanceErrorSet{Balance: balance.Balance{}, error: errors.New(balance.EmptyBalancesMessage)}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithNoDate(t *testing.T) {
	balances := balance.Balances{balance.Balance{}}
	expected := BalanceErrorSet{balance.Balance{}, nil}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithSingleDate(t *testing.T) {
	earliest, _ := balance.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), newTestMoney(t, 10, "GBP"))
	balances := balance.Balances{earliest}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest, _ := balance.New(date, newTestMoney(t, 10, "GBP"))
	other, _ := balance.New(date, newTestMoney(t, 20, "EUR"))
	balances := balance.Balances{earliest, other}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest, _ := balance.New(date1, newTestMoney(t, 10, "GBP"))
	other, _ := balance.New(date2, newTestMoney(t, 1, "GBP"))
	other2, _ := balance.New(date1, newTestMoney(t, 8765, "GBP"))
	other3, _ := balance.New(date3, newTestMoney(t, 489, "EUR"))
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
	latest, _ := balance.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC), newTestMoney(t, 10, "GBP"))
	balances := balance.Balances{latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	latest, _ := balance.New(date, newTestMoney(t, 10, "GBP"))
	other, _ := balance.New(date, newTestMoney(t, 20, "EUR"))
	balances := balance.Balances{other, latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	latest, _ := balance.New(date3, newTestMoney(t, 20, "GBP"))
	other1, _ := balance.New(date2, newTestMoney(t, 0, "GBP"))
	other2, _ := balance.New(date3, newTestMoney(t, 10, "EUR"))
	other3, _ := balance.New(date1, newTestMoney(t, 20, "YEN"))
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
	type testMoney struct {
		amount   int64
		currency string
	}
	testSets := []struct {
		moneys      []testMoney
		expectedSum money.Money
		error
	}{
		{},
		{
			moneys: []testMoney{
				{amount: 1, currency: "GBP"},
			},
			expectedSum: newTestMoney(t, 1, "GBP"),
		},
		{
			moneys: []testMoney{
				{amount: 1, currency: "GBP"},
				{amount: 2, currency: "GBP"},
			},
			expectedSum: newTestMoney(t, 3, "GBP"),
		},
		{
			moneys: []testMoney{
				{amount: 1, currency: "GBP"},
				{amount: 2, currency: "GBP"},
				{amount: -3, currency: "GBP"},
			},
			expectedSum: newTestMoney(t, 0, "GBP"),
		},

		{
			moneys: []testMoney{
				{amount: 1, currency: "GBP"},
				{amount: 2, currency: "EUR"},
			},
			error: money.CurrencyMismatchError{
				A: *innermoney.GetCurrency("GBP"),
				B: *innermoney.GetCurrency("EUR"),
			},
		},
	}

	now := time.Now()

	for i, testSet := range testSets {
		var bs balance.Balances
		for _, tsm := range testSet.moneys {
			m := newTestMoney(t, tsm.amount, tsm.currency)
			b, err := balance.New(now, m)
			common.FatalIfErrorf(t, err, "[%d] creating Balance for testing", i)
			bs = append(bs, b)
		}
		actual, err := bs.Sum()
		assert.Equal(t, testSet.error, err, "summing Balances ", testSet.moneys)
		if err != nil || len(bs) == 0 {
			continue
		}
		equal, err := actual.Equal(testSet.expectedSum)
		common.FatalIfErrorf(t, err, "[%d] Equalling", i)
		assert.True(t, equal, "[%d] %+v", i, testSet.moneys)
	}
}

func TestBalance_MarshalJSON(t *testing.T) {
	a, _ := balance.New(time.Now(), newTestMoney(t, 8237, "YEN"))
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Error marshalling json for testing: %s", err)
	}
	var b struct {
		Date  time.Time
		Money money.Money
	}
	if err := json.Unmarshal(jsonBytes, &b); err != nil {
		t.Fatalf("Error unmarshalling data into object: %s", err)
	}
	if equal, _ := a.Money().Equal(b.Money); !equal {
		t.Errorf("Expected %+v but got %+v.\njson: %s", a.Money(), b.Money, jsonBytes)
	}
	if !a.Date().Equal(b.Date) {
		t.Errorf("Expected %+v but got %+v.\njson: %s", a.Date(), b.Date, jsonBytes)
	}
}

func TestBalance_JSONLoop(t *testing.T) {
	a, _ := balance.New(time.Now(), newTestMoney(t, 8237, "USD"))
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Error marshalling json for testing: %s", err)
	}
	var b balance.Balance
	if err := json.Unmarshal(jsonBytes, &b); err != nil {
		t.Fatalf("Error unmarshaling bytes into Balance: %s", err)
	}
	if !a.Equal(b) {
		t.Fatalf("Expected %v, but got %v", a, b)
	}
}

func newTestMoney(t *testing.T, amount int64, currency string) money.Money {
	a, err := money.New(amount, currency)
	common.FatalIfError(t, err, "Creating Balance for testing")
	return *a
}
