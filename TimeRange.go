package GOHMoney

import (
	"github.com/lib/pq"
	"time"
)

// TimeRange represents a range of time that can be open ended at none, either or both ends.
type TimeRange struct {
	Start pq.NullTime
	End   pq.NullTime
}

// Validate checks the fields of TimeRange to ensure that if either the Start or End time is present, that the End time isn't before the Start time and returns an error if it is.
func (timeRange TimeRange) Validate() error {
	if timeRange.Start.Valid && timeRange.End.Valid && timeRange.End.Time.Before(timeRange.Start.Time) {
		return DateClosedBeforeDateOpenedError
	}
	return nil
}

// Contains returns true if the TimeRange contains given time.
// Contains will always return true when both the Start time and End time are not Valid
// Contains returns true if the time is on or after the TimeRange's Start time and before the TimeRange's End time.
func (timeRange TimeRange) Contains(time time.Time) bool {
	if timeRange.Start.Valid && time.Before(timeRange.Start.Time) {
		return false
	}
	if timeRange.End.Valid && (time.Equal(timeRange.End.Time) || time.After(timeRange.End.Time)) {
		return false
	}
	return true
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

