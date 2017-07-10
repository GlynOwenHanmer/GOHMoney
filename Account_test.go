package GOHMoney

import (
	"testing"
	"github.com/lib/pq"
)

func Test_IsOpen(t *testing.T) {
	testSets := []struct {
		Account
		IsOpen bool
	}{
		{
			Account: Account{},
			IsOpen:  true,
		},
		{
			Account: Account{DateClosed:pq.NullTime{Valid:false}},
			IsOpen:  true,
		},
		{
			Account: Account{DateClosed:pq.NullTime{Valid:true}},
			IsOpen:  false,
		},
	}
	for _, testSet := range testSets {
		actual := testSet.Account.IsOpen()
		if actual != testSet.IsOpen {
			t.Errorf("Account IsOpen expected %t, got %t. Account: %v", testSet.IsOpen, actual, testSet.Account)
		}
	}
}
