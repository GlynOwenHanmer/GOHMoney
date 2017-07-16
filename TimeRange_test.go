package GOHMoney

import (
	"testing"
	"github.com/lib/pq"
	"time"
	"bytes"
	"fmt"
)

func Test_Validate(t *testing.T) {
	testSets := []struct{
		TimeRange
		error
	}{
		{
			TimeRange:TimeRange{
				Start: pq.NullTime{Valid: false},
				End:   pq.NullTime{Valid: false},
			},
			error: nil,
		},
		{
			TimeRange:TimeRange{
				Start: pq.NullTime{Valid: true},
				End:   pq.NullTime{Valid: false},
			},
			error: nil,
		},
		{
			TimeRange:TimeRange{
				Start: pq.NullTime{Valid: false},
				End:   pq.NullTime{Valid: true},
			},
			error: nil,
		},
		{
			TimeRange:TimeRange{
				Start: pq.NullTime{Valid: true},
				End:   pq.NullTime{Valid: true},
			},
			error: nil,
		},
		{
			TimeRange:TimeRange{
				Start:pq.NullTime{
					Valid:true,
					Time:time.Now().AddDate(-1,0,0),
				},
				End:pq.NullTime{
					Valid:true,
					Time:time.Now(),
				},
			},
			error: nil,
		},
		{
			TimeRange:TimeRange{
				Start:pq.NullTime{
					Valid:true,
					Time:time.Now(),
				},
				End:pq.NullTime{
					Valid:true,
					Time:time.Now().AddDate(-1,0,0),
				},
			},
			error: DateClosedBeforeDateOpenedError,
		},
	}
	for _, testSet := range testSets {
		err := testSet.TimeRange.Validate()
		if err != testSet.error {
			t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", testSet.error, err)
		}
	}
}

func Test_Contains(t *testing.T) {
	testStartTime := time.Now()
	openRange := TimeRange{}
	openEnded := TimeRange{
		Start:pq.NullTime{
			Valid: true,
			Time:  testStartTime,
		},
	}
	openStarted := TimeRange{
		End:pq.NullTime{
			Valid: true,
			Time:  testStartTime,
		},
	}
	closedEnds := TimeRange{
		Start:pq.NullTime{
			Valid: true,
			Time:  testStartTime,
		},
		End:pq.NullTime{
			Valid: true,
			Time:  testStartTime.AddDate(1,0,0),
		},
	}

	testSets := []struct{
		TimeRange
		time.Time
		contains bool
	}{
		{
			TimeRange:openRange,
			contains:true,
		},
		{
			TimeRange: openRange,
			Time:      testStartTime,
			contains:  true,
		},
		{
			TimeRange: openEnded,
			Time:     testStartTime.AddDate(-1,0,0),
			contains: false,
		},
		{
			TimeRange:openEnded,
			Time:     testStartTime,
			contains: true,
		},
		{
			TimeRange:openEnded,
			Time:     testStartTime.AddDate(1,0,0),
			contains: true,
		},
		{
			TimeRange: openStarted,
			Time:     testStartTime.AddDate(-1,0,0),
			contains: true,
		},
		{
			TimeRange:openStarted,
			Time:     testStartTime,
			contains: false,
		},
		{
			TimeRange:openStarted,
			Time:     testStartTime.AddDate(1,0,0),
			contains: false,
		},
		{
			TimeRange:closedEnds,
			Time:testStartTime.AddDate(-2,0,0),
			contains:false,
		},
		{
			TimeRange:closedEnds,
			Time:testStartTime,
			contains:true,
		},
		{
			TimeRange:closedEnds,
			Time:testStartTime.AddDate(0,6,0),
			contains:true,
		},
		{
			TimeRange:closedEnds,
			Time:testStartTime.AddDate(1,0,0),
			contains:false,
		},
		{
			TimeRange:closedEnds,
			Time:testStartTime.AddDate(2,0,0),
			contains:false,
		},
	}
	for _, testSet := range testSets {
		contains := testSet.TimeRange.Contains(testSet.Time)
		if contains != testSet.contains {
			var message bytes.Buffer
			fmt.Fprint(&message, `Unexpected Contains result.`)
			fmt.Fprintf(&message, "\nExpected Contains: %t\nActual Contains  : %t", testSet.contains, contains)
			fmt.Fprintf(&message, "\nTimeRange: %+v", testSet.TimeRange)
			fmt.Fprintf(&message, "\nTime: %+v", testSet.Time)
			t.Error(message.String())
		}
	}
}

func Test_Equal(t *testing.T) {
	testSets := []struct {
		a, b TimeRange
		equal bool
	}{
		{
			a:TimeRange{},
			b:TimeRange{},
			equal:true,
		},
		{
			a:TimeRange{
				Start:pq.NullTime{
					Valid:true,
				},
			},
			b:TimeRange{},
			equal:false,
		},
		{
			a:TimeRange{
				End:pq.NullTime{
					Valid:true,
				},
			},
			b:TimeRange{},
			equal:false,
		},
	}
	for _, testSet := range testSets {
		if equal := testSet.a.Equal(testSet.b); equal != testSet.equal {
			t.Errorf(`Unexpected equal result.\nExpected: %t, Actual  : %t`, testSet.equal, equal)
		}
	}
}