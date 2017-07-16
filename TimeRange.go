package GOHMoney

import (
	"github.com/lib/pq"
)

// TimeRange represents a range of time that can be open ended at none, either or both ends.
type TimeRange struct {
	Start pq.NullTime
	End   pq.NullTime
}

// Validate checks the fields of TimeRange to ensure that if either the Start or End time is present, that the End time isn't before the Start time and returns an error if it is.
func (tr TimeRange) Validate() error {
	switch {
	case tr.Start.Valid == false,
		tr.End.Valid == false:
		return nil
	case tr.End.Time.Before(tr.Start.Time):
		return DateClosedBeforeDateOpenedError
	}
	return nil
}

// timeRangeValidationError holds an error describing an issue with a TimeRange object
type timeRangeValidationError string

// Error ensures that timeRangeValidationError adheres to the error interface
func (trve timeRangeValidationError) Error() string {
	return string(trve)
}

const (
	DateClosedBeforeDateOpenedError = timeRangeValidationError("Closed date is before opened date.")
)

