package account_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/GlynOwenHanmer/GOHMoney"
	"github.com/GlynOwenHanmer/GOHMoney/account"
)

func TestAccount_MarshalJSON(t *testing.T) {
	now := time.Now()
	testSets := []struct {
		start time.Time
		end   GOHMoney.NullTime
	}{
		{
			start: now,
			end:   GOHMoney.NullTime{},
		},
		{
			start: now,
			end: GOHMoney.NullTime{
				Valid: true,
				Time:  now.AddDate(1, 0, 0),
			},
		},
	}
	for _, set := range testSets {
		a, err := account.New("TEST ACCOUNT", set.start, set.end)
		if err != nil {
			t.Fatalf("Error creating new a for testings. Error: %s", err.Error())
		}
		bytes, err := json.Marshal(&a)
		if err != nil {
			t.Fatalf("Error marshalling json for testing. Error: %s", err.Error())
		}
		var unmarshalled account.Account
		err = json.Unmarshal(bytes, &unmarshalled)
		if err != nil {
			t.Errorf("Error unmarshalling Account json blob. Error: %s\njson: %s", err.Error(), bytes)
		}
		if unmarshalled.Name != a.Name {
			t.Errorf(`Unexpected name. Expected "%s" but got "%s"`, a.Name, unmarshalled.Name)
		}
		if !a.Start().Equal(unmarshalled.Start()) {
			t.Errorf("Unexpected a Start.\n\tExpected: %s\n\tActual  : %s", a.Start(), unmarshalled.Start())
		}
		if a.End().Valid != unmarshalled.End().Valid || !a.End().Time.Equal(unmarshalled.End().Time) {
			t.Errorf("Unexpected a End. \n\tExpected: %v\n\tActual  : %v", a.End(), unmarshalled.End())
		}
	}
}

func TestAccount_Equal(t *testing.T) {
	now := time.Now()
	a, err := account.New("A", now, GOHMoney.NullTime{})
	if err != nil {
		t.Errorf("Error creating account for testing: %s", err)
	}
	tests := []struct {
		name  string
		start time.Time
		end   GOHMoney.NullTime
		equal bool
	}{
		{"A", now, GOHMoney.NullTime{}, true},
		{"B", now, GOHMoney.NullTime{}, false},
		{"A", now.AddDate(-1, 0, 0), GOHMoney.NullTime{}, false},
		{"A", now, GOHMoney.NullTime{Valid: true, Time: now.Add(1)}, false},
		{"A", now.AddDate(-1, 0, 0), GOHMoney.NullTime{Valid: true, Time: now.Add(1)}, false},
		{"B", now.AddDate(-1, 0, 0), GOHMoney.NullTime{Valid: true, Time: now.Add(1)}, false},
	}
	for _, test := range tests {
		b, err := account.New(test.name, test.start, test.end)
		if err != nil {
			t.Errorf("Error creating account for testing: %s", err)
		}
		equal := a.Equal(b)
		if equal != test.equal {
			t.Errorf("Expected %t, but got %t.\nA: %v\nB: %v", test.equal, equal, a, b)
		}
	}
}
