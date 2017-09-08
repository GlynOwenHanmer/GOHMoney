package GOHMoney_test

import (
	"testing"
	"github.com/GlynOwenHanmer/GOHMoney"
	"time"
	"encoding/json"
	"github.com/lib/pq"
)

func TestAccount_MarshalJSON(t *testing.T) {
	now := time.Now()
	testSets := []struct{
		start time.Time
		end pq.NullTime
	}{
		{
			start:now,
			end:pq.NullTime{},
		},
		{
			start:now,
			end:pq.NullTime{
				Valid:true,
				Time:now.AddDate(1,0,0),
			},
		},
	}
	for _, set := range testSets {
		account, err := GOHMoney.NewAccount("TEST ACCOUNT",set.start,set.end)
		if err != nil {
			t.Fatalf("Error creating new account for testings. Error: %s", err.Error())
		}
		bytes, err := json.Marshal(&account)
		if err != nil {
			t.Fatalf("Error marshalling json for testing. Error: %s", err.Error())
		}
		var unmarshalled GOHMoney.Account
		err = json.Unmarshal(bytes, &unmarshalled)
		if err != nil {
			t.Errorf("Error unmarshalling Account json blob. Error: %s\njson: %s", err.Error(), bytes)
		}
		if unmarshalled.Name != account.Name {
			t.Errorf(`Unexpected name. Expected "%s" but got "%s"`, account.Name, unmarshalled.Name)
		}
		if !account.Start().Equal(unmarshalled.Start()) {
			t.Errorf("Unexpected account Start.\n\tExpected: %s\n\tActual  : %s", account.Start(), unmarshalled.Start())
		}
		if account.End().Valid != unmarshalled.End().Valid || !account.End().Time.Equal(unmarshalled.End().Time) {
			t.Errorf("Unexpected account End. \n\tExpected: %s\n\tActual  : %s", account.End(), unmarshalled.End())
		}	}
}

func TestAccount_Equal(t *testing.T) {
	now := time.Now()
	a, err := GOHMoney.NewAccount("A", now, pq.NullTime{})
	if err != nil {
		t.Errorf("Error creating account for testing: %s", err)
	}
	tests := []struct{
		name string
		start time.Time
		end pq.NullTime
		equal bool
	}{
		{"A", now, pq.NullTime{}, true},
		{"B", now, pq.NullTime{}, false},
		{"A", now.AddDate(-1,0,0), pq.NullTime{}, false},
		{"A", now, pq.NullTime{Valid:true, Time:now.Add(1)}, false},
		{"A", now.AddDate(-1,0,0), pq.NullTime{Valid:true, Time:now.Add(1)}, false},
		{"B", now.AddDate(-1,0,0), pq.NullTime{Valid:true, Time:now.Add(1)}, false},
	}
	for _, test := range tests {
		b, err := GOHMoney.NewAccount(test.name, test.start, test.end)
		if err != nil {
			t.Errorf("Error creating account for testing: %s", err)
		}
		equal := a.Equal(&b)
		if equal != test.equal {
			t.Errorf("Expected %s, but got %t.\nA: %s\nB: %s", test.equal, equal, a, b)
		}
	}
}