package GOHMoney

import (
	"testing"
	"github.com/lib/pq"
	"time"
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