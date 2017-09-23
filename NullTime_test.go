package GOHMoney

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func Test_NullTimeEqual(t *testing.T) {
	timeNow := time.Now()
	testSets := []struct {
		a, b  NullTime
		equal bool
	}{
		// a not Valid, b not Valid
		{
			a:     NullTime{},
			b:     NullTime{},
			equal: true,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b:     NullTime{},
			equal: false,
		},
		{
			a: NullTime{},
			b: NullTime{
				Time: timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b: NullTime{
				Time: timeNow,
			},
			equal: true,
		},

		// a Valid, b not Valid
		{
			a: NullTime{
				Valid: true,
			},
			b:     NullTime{},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b:     NullTime{},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
			},
			b: NullTime{
				Time: timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b: NullTime{
				Time: timeNow,
			},
			equal: false,
		},

		// a not Valid, b Valid
		{
			a: NullTime{},
			b: NullTime{
				Valid: true,
			},
			equal: false,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b: NullTime{
				Valid: true,
			},
			equal: false,
		},
		{
			a: NullTime{},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Time: timeNow,
			},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: false,
		},

		// a Valid, b Valid
		{
			a: NullTime{
				Valid: true,
			},
			b: NullTime{
				Valid: true,
			},
			equal: true,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b: NullTime{
				Valid: true,
			},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
			},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: false,
		},
		{
			a: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			b: NullTime{
				Valid: true,
				Time:  timeNow,
			},
			equal: true,
		},
	}

	for _, testSet := range testSets {
		if equal := testSet.a.equal(testSet.b); testSet.equal != equal {
			var message bytes.Buffer
			fmt.Fprintf(&message, "Unexpected equal result.\nExpected: %t\nActual  : %t", testSet.equal, equal)
			fmt.Fprintf(&message, "\na: %+v", testSet.a)
			fmt.Fprintf(&message, "\nb: %+v", testSet.b)
			t.Errorf(message.String())
		}
	}
}
