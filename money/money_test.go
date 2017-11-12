package money_test

import (
	"testing"

	"encoding/json"

	"github.com/glynternet/go-money/money"
	"github.com/glynternet/go-money/currency"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c, err := currency.NewCode("RIN")
	assert.Nil(t, err)
	m := money.New(123, *c)
	assert.NotNil(t, m)
	assert.Equal(t, "RIN", m.Currency().String())
	assert.Equal(t, 123, m.Amount())
}

func TestJSON(t *testing.T) {
	c, err := currency.NewCode("RIN")
	assert.Nil(t, err)
	ma := money.New(9876, *c)
	bs, err := json.Marshal(ma)
	assert.Nil(t, err)
	mb, err := money.UnmarshalJSON(bs)
	assert.Nil(t, err, string(bs))
	assert.Equal(t, ma, *mb)
}

//func TestMoneyEqual(t *testing.T) {
//	testSets := []struct {
//		a, b  money.Money
//		equal bool
//		error
//	}{
//		{
//			equal: false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     *newMoneyIgnoreError(0, ""),
//			equal: false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     money.Money{},
//			equal: false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     *newMoneyIgnoreError(0, "EUR"),
//			equal: false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     *newMoneyIgnoreError(-10, "GBP"),
//			equal: false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			b:     *newMoneyIgnoreError(1023, "GBP"),
//			equal: false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     *newMoneyIgnoreError(103, "GBP"),
//			b:     *newMoneyIgnoreError(1023, "GBP"),
//			equal: false,
//		},
//		{
//			a:     *newMoneyIgnoreError(1023, "USD"),
//			b:     *newMoneyIgnoreError(1023, "USD"),
//			equal: true,
//		},
//	}
//	for i, ts := range testSets {
//		equal, err := ts.a.Equal(ts.b)
//		assert.Equal(t, ts.error, err, "[%+v] a: %v, b: %v", i, ts.a, ts.b)
//		assert.Equal(t, ts.equal, equal, "[%d] a: %+v, b: %+v", i, ts.a, ts.b)
//	}
//}
//
//func TestMoneyAdd(t *testing.T) {
//	testSets := []struct {
//		a, b, sum money.Money
//		error
//	}{
//		{
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     money.Money{},
//			b:     money.Money{},
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     *newMoneyIgnoreError(1, "EUR"),
//			b:     money.Money{},
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:     *newMoneyIgnoreError(1, "EUR"),
//			b:     *newMoneyIgnoreError(2, "GBP"),
//			error: money.CurrencyMismatchError{A: *money2.GetCurrency("EUR"), B: *money2.GetCurrency("GBP")},
//		},
//		{
//			a:     *newMoneyIgnoreError(-3, "USD"),
//			b:     *newMoneyIgnoreError(-10, "USD"),
//			sum:   *newMoneyIgnoreError(-13, "USD"),
//			error: nil,
//		},
//	}
//	for i, ts := range testSets {
//		sum, err := ts.a.Add(ts.b)
//		assert.Equal(t, ts.error, err, "[%d] a: %v, b: %v", i, ts.error, err)
//		if err != nil {
//			continue
//		}
//		equal, _ := sum.Equal(ts.sum)
//		assert.True(t, equal, "[%d] a: %v, b: %v, sum: %v", i, ts.a, ts.b, sum)
//	}
//}
//
//func TestMoneyJSONLoop(t *testing.T) {
//	a, err := money.NewCode(934, "YEN")
//	common.FatalIfError(t, err, "Creating Money")
//	jsonBytes, err := json.Marshal(a)
//	if err != nil {
//		t.Fatalf("Error marshalling json for testing: %s", err)
//	}
//	var b money.Money
//	if err := json.Unmarshal(jsonBytes, &b); err != nil {
//		t.Fatalf("Error unmarshaling bytes into object: %s", err)
//	}
//	if equal, _ := a.Equal(b); !equal {
//		t.Fatalf("Expected %v, but got %v", a, b)
//	}
//}
//
//func TestMoneySameCurrency(t *testing.T) {
//	same, err := (*newMoneyIgnoreError(234, "EUR")).SameCurrency()
//	assert.True(t, same)
//	assert.Nil(t, err)
//
//	same, err = money.Money{}.SameCurrency()
//	assert.True(t, same)
//	assert.Nil(t, err)
//
//	testSets := []struct {
//		a, b money.Money
//		bool
//		error
//	}{
//		{
//			bool:  false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			b:     *newMoneyIgnoreError(123, "EUR"),
//			bool:  false,
//			error: money.ErrNoCurrency,
//		},
//		{
//			a:    *newMoneyIgnoreError(123, "GBP"),
//			b:    *newMoneyIgnoreError(123, "EUR"),
//			bool: false,
//		},
//		{
//			a:    *newMoneyIgnoreError(123, "GBP"),
//			b:    *newMoneyIgnoreError(987, "GBP"),
//			bool: true,
//		},
//	}
//	for i, ts := range testSets {
//		same, err := ts.a.SameCurrency(ts.b)
//		assert.Equal(t, ts.bool, same, "[%d] a: %+v, b: %+v", i, ts.a, ts.b)
//		assert.Equal(t, ts.error, err, "[%d] a: %+v, b: %+v", i, ts.a, ts.b)
//	}
//}
