package money_test

import (
	"testing"

	"encoding/json"

	"github.com/GlynOwenHanmer/GOHMoney/money"
	money2 "github.com/rhymond/go-money"
)

func TestMoneyCurrency(t *testing.T) {
	testSets := []struct {
		money.Money
		code string
	}{
		{
			code: "GBP",
		},
		{
			Money: money.New(0),
			code:  "GBP",
		},
	}
	for _, ts := range testSets {
		expected := money2.GetCurrency(ts.code)
		actual := ts.Money.Currency()
		if actual != *expected {
			t.Errorf("Expected %v but got %v", expected, actual)
		}
	}
}

func TestMoneyAmount(t *testing.T) {
	testSets := []struct {
		money.Money
		amount int64
	}{
		{
			amount: -99,
			Money:  money.New(-99),
		},
		{
			Money: money.New(0),
		},
		{
			amount: 9876,
			Money:  money.New(9876),
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
	}{
		{
			equal: true,
		},
		{
			a:     money.New(0),
			equal: true,
		},
		{
			b:     money.New(0),
			equal: true,
		},

		{
			a:     money.New(-10),
			equal: false,
		},
		{
			b:     money.New(1023),
			equal: false,
		},
	}
	for _, ts := range testSets {
		equal, _ := ts.a.Equal(ts.b)
		if equal != ts.equal {
			t.Errorf("Expected %t but got %t", ts.equal, equal)
		}
	}
}

func TestMoneyAdd(t *testing.T) {
	testSets := []struct {
		a, b, sum money.Money
	}{
		{
			sum: money.New(0),
		},
		{
			a: money.Money{},
		},
	}
	for _, ts := range testSets {
		sum, _ := ts.a.Add(ts.b)
		if equal, _ := sum.Equal(ts.sum); !equal {
			t.Errorf("Expected %v, got %v", ts.sum, sum)
		}
	}
}

func TestMoneyJSONLoop(t *testing.T) {
	a := money.New(934)
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
