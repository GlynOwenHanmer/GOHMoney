package account

import (
	"time"

	gtime "github.com/glynternet/go-time"
)

// Option is a function that takes a pointer to an Account returning an error.
// The idea of Option is to alter a Account object
type Option func(*Account) error

// CloseTime returns an Option that will set the close time on an Account object.
func CloseTime(t time.Time) Option {
	if t.IsZero() {
		return nil
	}
	return func(a *Account) error {
		a.timeRange.End = gtime.NullTime{Valid: true, Time: t}
		return nil
	}
}