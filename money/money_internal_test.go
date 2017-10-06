package money

import (
	"testing"
)

func TestMoneyAdd(t *testing.T) {
	testSets := []struct {
		a, b, sum Money
	}{
		{
			sum: GBP(0),
		},
	}
	for _, ts := range testSets {
		sum, _ := ts.a.Add(ts.b)
		if equal, _ := sum.Equal(ts.sum); !equal {
			t.Errorf("Expected %v, got %v", ts.sum, sum)
		}
	}
}
