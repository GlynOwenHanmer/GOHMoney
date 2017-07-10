package GOHMoney

import (
	"testing"
	"errors"
	"time"
)

type BalanceErrorSet struct {
	Balance
	error
}

func Test_Earliest_EmptyBalances(t *testing.T) {
	balances := Balances{}
	expected := BalanceErrorSet{Balance: Balance{}, error: errors.New(EmptyBalancesMessage)}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithNoDate(t *testing.T){
	balances := Balances{Balance{}}
	expected := BalanceErrorSet{Balance{},nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithSingleDate(t *testing.T){
	earliest := Balance{Date:time.Date(2000,1,1,1,1,1,1,time.UTC),Amount:10}
	balances := Balances{earliest}
	expected := BalanceErrorSet{earliest,nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithSameDate(t *testing.T){
	date := time.Date(2000,1,1,1,1,1,1,time.UTC)
	earliest := Balance{Date:date,Amount:10}
	balances := Balances{earliest, Balance{Date:date,Amount:20}}
	expected := BalanceErrorSet{earliest,nil}
	testEarliestSet(t, expected, balances)
}

func Test_Earliest_BalancesWithMultipleDates(t *testing.T){
	date1 := time.Date(2000,1,1,1,1,1,1,time.UTC)
	date2 := time.Date(2001,1,1,1,1,1,1,time.UTC)
	date3 := time.Date(2002,1,1,1,1,1,1,time.UTC)
	earliest := Balance{Date: date1,Amount: 10}
	balances := Balances{Balance{Date:date2}, earliest, Balance{Date: date1,Amount: 20}, Balance{Date:date3}}
	expected := BalanceErrorSet{earliest,nil}
	testEarliestSet(t, expected, balances)
}

func testEarliestSet(t *testing.T, expected BalanceErrorSet, balances Balances) {
	actualBalance, actualError := balances.Earliest()
	actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
	testBalanceResults(t, expected, actual)
}

func Test_Latest_EmptyBalances(t *testing.T) {
	balances := Balances{}
	expected := BalanceErrorSet{Balance: Balance{}, error: errors.New(EmptyBalancesMessage)}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithNoDate(t *testing.T){
	balances := Balances{Balance{}}
	expected := BalanceErrorSet{Balance{},nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSingleDate(t *testing.T){
	latest := Balance{Date: time.Date(2000,1,1,1,1,1,1,time.UTC),Amount: 10}
	balances := Balances{latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithSameDate(t *testing.T){
	date := time.Date(2000,1,1,1,1,1,1,time.UTC)
	latest := Balance{Date: date,Amount: 10}
	balances := Balances{Balance{Date: date,Amount: 20}, latest}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func Test_Latest_BalancesWithMultipleDates(t *testing.T){
	date1 := time.Date(2000,1,1,1,1,1,1,time.UTC)
	date2 := time.Date(2001,1,1,1,1,1,1,time.UTC)
	date3 := time.Date(2002,1,1,1,1,1,1,time.UTC)
	latest := Balance{Date: date3,Amount: 10}
	balances := Balances{Balance{Date:date2}, Balance{Date: date3}, latest, Balance{Date: date1,Amount: 20}}
	expected := BalanceErrorSet{latest, nil}
	testLatestSet(t, expected, balances)
}

func testLatestSet(t *testing.T, expected BalanceErrorSet, balances Balances) {
	actualBalance, actualError := balances.Latest()
	actual := BalanceErrorSet{Balance: actualBalance, error: actualError}
	testBalanceResults(t, expected, actual)
}

func testBalanceResults(t *testing.T, expected BalanceErrorSet, actual BalanceErrorSet) {
	if expected.error != actual.error {
		switch {
		case expected.error == nil:
			t.Errorf("Error unexpected\nExpected: %s\nActual  : %s", nil, actual)
		case actual.error == nil:
			t.Errorf("Error unexpected\nExpected: %s\nActual  : %s", expected, nil)
		case expected.error.Error() == actual.error.Error():
			break
		default:
			t.Errorf("Error unexpected\nExpected: %s\nActual  : %s", expected, actual)
		}
	}
	if expected.Balance != actual.Balance {
		t.Errorf("Balance unexpected\nExpected: %s\nActual  : %s", expected.Balance, actual.Balance)
	}
}

func TestBalances_Sum(t *testing.T) {
	testSets := []struct{
		Balances
		expectedSum float32
	}{
		{
			expectedSum:0,
		},
		{
			Balances: Balances{
				Balance{Amount:1},
			},
			expectedSum:1,
		},
		{
			Balances: Balances{
				Balance{Amount:1},
				Balance{Amount:2},
			},
			expectedSum:3,
		},

		{
			Balances: Balances{
				Balance{Amount:1},
				Balance{Amount:2},
				Balance{Amount:-3},
			},
			expectedSum:0,
		},
	}

	for _, testSet := range testSets {
		actual := testSet.Balances.Sum()
		if testSet.expectedSum != actual {
			t.Errorf("Unexpected sum.\nExpected: %f\nActual  : %f\nBalances: %v", testSet.expectedSum, actual, testSet,Balances{})
		}
	}
}