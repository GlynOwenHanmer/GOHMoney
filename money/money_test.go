package money_test

import (
	"testing"

	"encoding/json"

	"github.com/GlynOwenHanmer/GOHMoney/money"
	money2 "github.com/rhymond/go-money"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func newMoneyIgnoreError(a int64, c string) *money.Money {
	m, _ := money.New(a, c)
	return m
}

func TestMoneyCurrency(t *testing.T) {
	testSets := []struct {
		money.Money
		code string
		error
	}{
		{
			error: money.ErrNoCurrency,
			code:  "",
		},
		{
			Money: *newMoneyIgnoreError(123, ""),
			error: money.ErrNoCurrency,
			code:  "",
		},
		{
			Money: *newMoneyIgnoreError(0, "GBP"),
			code:  "GBP",
		},
	}
	for i, ts := range testSets {
		expected := money2.GetCurrency(ts.code)
		if expected == nil {
			expected = new(money2.Currency)
			*expected = money2.Currency{}
		}
		actual, err := ts.Money.Currency()
		assert.Equal(t, ts.error, err, "[%d] Money: %+v", i, ts.Money)
		assert.Equal(t, *expected, actual, "[%d] Money: %+v", i, ts.Money)
	}
}

func TestMoneyAmount(t *testing.T) {
	testSets := []struct {
		money.Money
		amount int64
	}{
		{
			amount: -99,
			Money:  *newMoneyIgnoreError(-99, "GBP"),
		},
		{
			amount: -99,
			Money:  *newMoneyIgnoreError(-99, "EUR"),
		},
		{
			Money: *newMoneyIgnoreError(0, "EUR"),
		},
		{
			amount: 9876,
			Money:  *newMoneyIgnoreError(9876, "GBP"),
		},
	}
	for _, ts := range testSets {
		expected := ts.amount
		actual := ts.Money.Amount()
		if actual != expected {
			t.Errorf("Expected %v but got %v", expected, actual)
		}
	}
}

func TestMoneyEqual(t *testing.T) {
	testSets := []struct {
		a, b  money.Money
		equal bool
		error
	}{
		{
			equal: true,
			error: money.ErrNoCurrency,
		},
		{
			a:     *newMoneyIgnoreError(0, ""),
			equal: true,
			error: money.ErrNoCurrency,
		},
		{
			a:     *newMoneyIgnoreError(0, "EUR"),
			equal: false,
			error: money.ErrNoCurrency,
		},
		{
			b:     *newMoneyIgnoreError(0, "GBP"),
			equal: true,
		},

		{
			a:     *newMoneyIgnoreError(-10, "GBP"),
			equal: false,
		},
		{
			b:     money.GBP(1023),
			equal: false,
		},
	}
	for i, ts := range testSets {
		equal, err := ts.a.Equal(ts.b)
		assert.Equal(t, err, ts.error, "[%+v] a: %v, b: %v", i, ts.a, ts.b)
		assert.Equal(t, ts.equal, equal, "[%d] a: %+v, b: %+v", i, ts.a, ts.b)
	}
}

func TestMoneyAdd(t *testing.T) {
	testSets := []struct {
		a, b, sum money.Money
		error
	}{
		{
			error: money.ErrNoCurrency,
		},
		{
			a:     money.Money{},
			b:     money.Money{},
			sum:   money.Money{},
			error: money.ErrNoCurrency,
		},
		{
			a:     *newMoneyIgnoreError(1, "EUR"),
			b:     *newMoneyIgnoreError(2, "GBP"),
			error: money.CurrencyMismatchError{A: *money2.GetCurrency("EUR"), B: *money2.GetCurrency("GBP")},
		},
		{
			a:     *newMoneyIgnoreError(-3, "USD"),
			b:     *newMoneyIgnoreError(-10, "USD"),
			sum:   *newMoneyIgnoreError(-13, "USD"),
			error: nil,
		},
	}
	for _, ts := range testSets {
		sum, err := ts.a.Add(ts.b)
		if err != ts.error {
			t.Fatalf("Expected %+v, got %+v", ts.error, err)
		}
		if equal, _ := sum.Equal(ts.sum); !equal {
			t.Errorf("Expected %v, got %v", ts.sum, sum)
		}
	}
}

func TestMoneyJSONLoop(t *testing.T) {
	a := money.GBP(934)
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Error marshalling json for testing: %s", err)
	}

	var b money.Money
	if err := json.Unmarshal(jsonBytes, &b); err != nil {
		t.Fatalf("Error unmarshaling bytes into object: %s", err)
	}
	if equal, _ := a.Equal(b); !equal {
		t.Fatalf("Expected %v, but got %v", a, b)
	}
}

func TestMoneySameCurrency(t *testing.T) {
	same, err := (*newMoneyIgnoreError(234,"EUR")).SameCurrency()
	assert.True(t, same)
	assert.Nil(t, err)

	same, err = money.Money{}.SameCurrency()
	assert.True(t, same)
	assert.Nil(t, err)

	testSets := []struct {
		a, b money.Money
		bool
		error
	}{
		{
			bool:  false,
			error: money.ErrNoCurrency,
		},
		{
			b:    *newMoneyIgnoreError(123, "EUR"),
			bool: false,
			error: money.ErrNoCurrency,
		},
		{
			a:    *newMoneyIgnoreError(123, "GBP"),
			b:    *newMoneyIgnoreError(123, "EUR"),
			bool: false,
			error: money.CurrencyMismatchError{A:money2.Currency{Code:"GBP"},B:money2.Currency{Code:"EUR"}},
		},
		{
			a:    *newMoneyIgnoreError(123, "GBP"),
			b:    *newMoneyIgnoreError(987, "GBP"),
			bool: true,
		},
	}
	for i, ts := range testSets {
		fmt.Printf("Test Set - %d\n", i)
		same, err := ts.a.SameCurrency(ts.b)
		assert.Equal(t, ts.bool, same, "[%d] a: %+v, b: %+v", i, ts.a, ts.b)
		assert.Equal(t, ts.error, err, "[%d] a: %+v, b: %+v", i, ts.a, ts.b)
	}
}