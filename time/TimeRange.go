package time

import (
	"time"
)

// Range represents a range of time that can be open ended at none, either or both ends.
type Range struct {
	Start NullTime
	End   NullTime
}

// Equal returns true if two Range objects have matching Start and End NullTimes
func (r Range) Equal(r2 Range) bool {
	if !NullTime(r.Start).Equal(NullTime(r2.Start)) || !NullTime(r.End).Equal(NullTime(r2.End)) {
		return false
	}
	return true
}

// Validate checks the fields of Range to ensure that if either the Start or End time is present, that the End time isn't before the Start time and returns an error if it is.
func (r Range) Validate() error {
	if r.Start.Valid && r.End.Valid && r.End.Time.Before(r.Start.Time) {
		return EndTimeBeforeStartTime
	}
	return nil
}

// Contains returns true if the Range contains given time.
// Contains will always return true when both the Start time and End time are not Valid
// Contains returns true if the time is on or after the Range's Start time and before the Range's End time.
func (r Range) Contains(time time.Time) bool {
	if r.Start.Valid && time.Before(r.Start.Time) {
		return false
	}
	if r.End.Valid && !r.End.Time.After(time) {
		return false
	}
	return true
}

// rangeValidationError holds an error describing an issue with a Range object
type rangeValidationError string

// Error ensures that rangeValidationError adheres to the error interface
func (err rangeValidationError) Error() string {
	return string(err)
}

const (
	// EndTimeBeforeStartTime is a rangeValidationError that is returned when a Balance's closing date is set to before the opening date.
	EndTimeBeforeStartTime = rangeValidationError("end time is before start time.")
)
