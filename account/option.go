package account

import (
	"time"
	gtime "github.com/glynternet/go-time"
)

type Option func(*Account) error

func CloseTime(time time.Time) Option {
	if time.IsZero() {
		return nil
	}
	return func(a *Account) error {
		a.timeRange.End = gtime.NullTime{Valid:true, Time: time}
		return nil
	}
}
