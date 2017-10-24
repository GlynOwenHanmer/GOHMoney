package account

import (
	"time"

	gtime "github.com/glynternet/go-time"
)

type Option func(*Account) error

func CloseTime(t time.Time) Option {
	if t.IsZero() {
		return nil
	}
	return func(a *Account) error {
		a.timeRange.End = gtime.NullTime{Valid: true, Time: t}
		return nil
	}
}