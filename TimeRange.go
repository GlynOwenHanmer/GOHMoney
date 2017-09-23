package GOHMoney

import (
	"time"
)

// TimeRange represents a range of time that can be open ended at none, either or both ends.
type TimeRange struct {
	Start NullTime
	End   NullTime
}

// equal returns true if two TimeRange objects have matching Start and End NullTimes
func (tr TimeRange) Equal(tr2 TimeRange) bool {
	if !NullTime(tr.Start).Equal(NullTime(tr2.Start)) || !NullTime(tr.End).Equal(NullTime(tr2.End)) {
		return false
	}
	return true
}

// Validate checks the fields of TimeRange to ensure that if either the Start or End time is present, that the End time isn't before the Start time and returns an error if it is.
func (tr TimeRange) Validate() error {
	if tr.Start.Valid && tr.End.Valid && tr.End.Time.Before(tr.Start.Time) {
		return DateClosedBeforeDateOpenedError
	}
	return nil
}

// Contains returns true if the TimeRange contains given time.
// Contains will always return true when both the Start time and End time are not Valid
// Contains returns true if the time is on or after the TimeRange's Start time and before the TimeRange's End time.
func (tr TimeRange) Contains(time time.Time) bool {
	if tr.Start.Valid && time.Before(tr.Start.Time) {
		return false
	}
	if tr.End.Valid && !tr.End.Time.After(time) {
		return false
	}
	return true
}

// timeRangeValidationError holds an error describing an issue with a TimeRange object
type timeRangeValidationError string

// Error ensures that timeRangeValidationError adheres to the error interface
func (err timeRangeValidationError) Error() string {
	return string(err)
}

const (
	DateClosedBeforeDateOpenedError = timeRangeValidationError("Closed date is before opened date.")
)
