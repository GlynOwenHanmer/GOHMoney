package balance_test

import (
	"testing"
	"github.com/GlynOwenHanmer/GOHMoney/balance"
)

func TestNewMoney(t *testing.T) {
	amounts := []int64{-100,0,3000}
	for _, a := range amounts {
		m := balance.NewMoney(a)
		actual := m.Amount()
		if actual != a {
			t.Errorf("Expected %d, got %d", a, actual)
		}
	}
}

func TestMoney_Equals(t *testing.T) {
	a := balance.Money{}
	b := balance.Money{}
	a.Equals(b)
}