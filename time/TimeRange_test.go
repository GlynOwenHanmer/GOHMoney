package time

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func Test_Validate(t *testing.T) {
	testSets := []struct {
		Range
		error
	}{
		{
			Range: Range{
				Start: NullTime{Valid: false},
				End:   NullTime{Valid: false},
			},
			error: nil,
		},
		{
			Range: Range{
				Start: NullTime{Valid: true},
				End:   NullTime{Valid: false},
			},
			error: nil,
		},
		{
			Range: Range{
				Start: NullTime{Valid: false},
				End:   NullTime{Valid: true},
			},
			error: nil,
		},
		{
			Range: Range{
				Start: NullTime{Valid: true},
				End:   NullTime{Valid: true},
			},
			error: nil,
		},
		{
			Range: Range{
				Start: NullTime{
					Valid: true,
					Time:  time.Now().AddDate(-1, 0, 0),
				},
				End: NullTime{
					Valid: true,
					Time:  time.Now(),
				},
			},
			error: nil,
		},
		{
			Range: Range{
				Start: NullTime{
					Valid: true,
					Time:  time.Now(),
				},
				End: NullTime{
					Valid: true,
					Time:  time.Now().AddDate(-1, 0, 0),
				},
			},
			error: EndTimeBeforeStartTime,
		},
	}
	for _, testSet := range testSets {
		err := testSet.Range.Validate()
		if err != testSet.error {
			t.Errorf("Unexpected error.\nExpected: %s\nActual  : %s", testSet.error, err)
		}
	}
}

func Test_Contains(t *testing.T) {
	testStartTime := time.Now()
	openRange := Range{}
	openEnded := Range{
		Start: NullTime{
			Valid: true,
			Time:  testStartTime,
		},
	}
	openStarted := Range{
		End: NullTime{
			Valid: true,
			Time:  testStartTime,
		},
	}
	closedEnds := Range{
		Start: NullTime{
			Valid: true,
			Time:  testStartTime,
		},
		End: NullTime{
			Valid: true,
			Time:  testStartTime.AddDate(1, 0, 0),
		},
	}

	testSets := []struct {
		Range
		time.Time
		contains bool
	}{
		{
			Range:    openRange,
			contains: true,
		},
		{
			Range:    openRange,
			Time:     testStartTime,
			contains: true,
		},
		{
			Range:    openEnded,
			Time:     testStartTime.AddDate(-1, 0, 0),
			contains: false,
		},
		{
			Range:    openEnded,
			Time:     testStartTime,
			contains: true,
		},
		{
			Range:    openEnded,
			Time:     testStartTime.AddDate(1, 0, 0),
			contains: true,
		},
		{
			Range:    openStarted,
			Time:     testStartTime.AddDate(-1, 0, 0),
			contains: true,
		},
		{
			Range:    openStarted,
			Time:     testStartTime,
			contains: false,
		},
		{
			Range:    openStarted,
			Time:     testStartTime.AddDate(1, 0, 0),
			contains: false,
		},
		{
			Range:    closedEnds,
			Time:     testStartTime.AddDate(-2, 0, 0),
			contains: false,
		},
		{
			Range:    closedEnds,
			Time:     testStartTime,
			contains: true,
		},
		{
			Range:    closedEnds,
			Time:     testStartTime.AddDate(0, 6, 0),
			contains: true,
		},
		{
			Range:    closedEnds,
			Time:     testStartTime.AddDate(1, 0, 0),
			contains: false,
		},
		{
			Range:    closedEnds,
			Time:     testStartTime.AddDate(2, 0, 0),
			contains: false,
		},
	}
	for _, testSet := range testSets {
		contains := testSet.Range.Contains(testSet.Time)
		if contains != testSet.contains {
			var message bytes.Buffer
			fmt.Fprint(&message, `Unexpected Contains result.`)
			fmt.Fprintf(&message, "\nExpected Contains: %t\nActual Contains  : %t", testSet.contains, contains)
			fmt.Fprintf(&message, "\nTimeRange: %+v", testSet.Range)
			fmt.Fprintf(&message, "\nTime: %+v", testSet.Time)
			t.Error(message.String())
		}
	}
}

func Test_Equal(t *testing.T) {
	testSets := []struct {
		a, b  Range
		equal bool
	}{
		{
			a:     Range{},
			b:     Range{},
			equal: true,
		},
		{
			a: Range{
				Start: NullTime{
					Valid: true,
				},
			},
			b:     Range{},
			equal: false,
		},
		{
			a: Range{
				End: NullTime{
					Valid: true,
				},
			},
			b:     Range{},
			equal: false,
		},
	}
	for _, testSet := range testSets {
		if equal := testSet.a.Equal(testSet.b); equal != testSet.equal {
			t.Errorf(`Unexpected Equal result.\nExpected: %t, Actual  : %t`, testSet.equal, equal)
		}
	}
}
