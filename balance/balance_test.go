package balance_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"fmt"

	"github.com/glynternet/GOHMoney/balance"
	"github.com/glynternet/GOHMoney/common"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var expected error = balance.ZeroDate
	invalidTime := time.Time{}
	_, err := balance.New(invalidTime)
	assert.Equal(t, expected, err)

	expected = nil
	validTime := time.Now()
	_, err = balance.New(validTime, balance.Amount(11))
	assert.Equal(t, expected, err)
}

func TestBalance_Equal(t *testing.T) {
	newBalance := func(t *testing.T, date time.Time, amount int64) balance.Balance {
		b, err := balance.New(date, balance.Amount(amount))
		common.FatalIfError(t, err, "Creating new balance")
		return b
	}

	now := time.Now()
	a := newBalance(t, now, 123)
	b := newBalance(t, now, 123)
	assert.True(t, a.Equal(b))

	b = newBalance(t, now, -123)
	assert.True(t, !a.Equal(b))

	b = newBalance(t, now, 123)
	assert.True(t, a.Equal(b))

	b = newBalance(t, now.Add(1), 123)
	assert.True(t, !a.Equal(b))
}

type BalanceErrorSet struct {
	balance.Balance
	error
}

func TestBalances_Earliest_EmptyBalances(t *testing.T) {
	balances := balance.Balances{}
	expected := BalanceErrorSet{error: errors.New(balance.EmptyBalancesMessage)}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithSingleDate(t *testing.T) {
	earliest, _ := balance.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC))
	balances := balance.Balances{earliest}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest, _ := balance.New(date, balance.Amount(10))
	other, _ := balance.New(date, balance.Amount(20))
	balances := balance.Balances{earliest, other}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func TestBalances_Earliest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	earliest, _ := balance.New(date1, balance.Amount(10))
	other, _ := balance.New(date2, balance.Amount(1))
	other2, _ := balance.New(date1, balance.Amount(8237))
	other3, _ := balance.New(date3, balance.Amount(489))
	balances := balance.Balances{other, earliest, other2, other3}
	expected := BalanceErrorSet{earliest, nil}
	testEarliestSet(t, expected, balances)
}

func testEarliestSet(t *testing.T, expected BalanceErrorSet, balances balance.Balances) {
	actualBalance, actualError := balances.Earliest()
	actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
	res := testBalanceResults(t, expected, actual)
	if len(res) > 0 {
		t.Errorf("%s. Balances: %+v", res, balances)
	}
}

func Test_Latest_EmptyBalances(t *testing.T) {
	balances := balance.Balances{}
	expected := BalanceErrorSet{error: errors.New(balance.EmptyBalancesMessage)}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSingleDate(t *testing.T) {
	latest, _ := balance.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC))
	balances := balance.Balances{latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSameDate(t *testing.T) {
	date := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	latest, _ := balance.New(date, balance.Amount(10))
	other, _ := balance.New(date, balance.Amount(20))
	balances := balance.Balances{other, latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithMultipleDates(t *testing.T) {
	date1 := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	date2 := time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC)
	date3 := time.Date(2002, 1, 1, 1, 1, 1, 1, time.UTC)
	latest, _ := balance.New(date3, balance.Amount(20))
	other1, _ := balance.New(date2)
	other2, _ := balance.New(date3, balance.Amount(10))
	other3, _ := balance.New(date1, balance.Amount(-20))
	balances := balance.Balances{other1, other2, latest, other3}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func testLatestSet(t *testing.T, expected BalanceErrorSet, balances balance.Balances) {
	actualBalance, actualError := balances.Latest()
	actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
	res := testBalanceResults(t, expected, actual)
	if len(res) > 0 {
		t.Errorf("%s. Balances: %+v", res, balances)
	}
}

func testBalanceResults(t *testing.T, expected BalanceErrorSet, actual BalanceErrorSet) (message string) {
	if expected.error != actual.error {
		switch {
		case expected.error == nil:
			message = fmt.Sprintf("Expected no error but got %v", actual)
		case actual.error == nil:
			message = fmt.Sprintf("Error error (%v) but didn't get one", expected)
		case expected.error.Error() == actual.error.Error():
			break
		default:
			message = fmt.Sprintf("Error unexpected\nExpected: %s\nActual  : %s", expected, actual)
		}
	}
	assert.Equal(t, expected.Balance, actual.Balance)
	return
}

func TestBalances_Sum(t *testing.T) {
	testSets := []struct {
		amounts []int64
		sum     int64
	}{
		{},
		{
			amounts: []int64{1},
			sum:     1,
		},
		{
			amounts: []int64{1, 2},
			sum:     3,
		},
		{
			amounts: []int64{1, 2, -3},
			sum:     0,
		},
	}

	now := time.Now()

	for i, testSet := range testSets {
		var bs balance.Balances
		for _, tsm := range testSet.amounts {
			b, err := balance.New(now, balance.Amount(tsm))
			common.FatalIfErrorf(t, err, "[%d] creating balance for testing", i)
			bs = append(bs, b)
		}
		assert.Equal(t, testSet.sum, bs.Sum())
	}
}

func TestBalance_MarshalJSON(t *testing.T) {
	a, err := balance.New(time.Now(), balance.Amount(921368))
	common.FatalIfError(t, err, "Creating Balance")
	jsonBytes, err := json.Marshal(a)
	common.FatalIfError(t, err, "Marshalling JSON")

	var b struct {
		Date   time.Time
		Amount int64
	}
	err = json.Unmarshal(jsonBytes, &b)
	common.FatalIfError(t, err, "Unmarshalling data")
	assert.True(t, a.Date().Equal(b.Date), "json: %s", jsonBytes)
	assert.Equal(t, a.Amount(), b.Amount, "json: %s", jsonBytes)
}

func TestBalance_JSONLoop(t *testing.T) {
	a, _ := balance.New(time.Now(), balance.Amount(8237))
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Error marshalling json for testing: %s", err)
	}
	var b balance.Balance
	if err := json.Unmarshal(jsonBytes, &b); err != nil {
		t.Fatalf("Error unmarshaling bytes into balance: %s", err)
	}
	if !a.Equal(b) {
		t.Fatalf("Expected %v, but got %v", a, b)
	}
}
