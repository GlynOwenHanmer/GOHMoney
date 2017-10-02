package balance

import (
	"fmt"
	"time"

	gohtime "github.com/GlynOwenHanmer/go-time"
)

// DateOutOfAccountTimeRange is a type returned when the date of a Balance is not contained within the Range of the Account that holds it.
// BalanceDate and AccountTimeRange fields are present and provide the exact detail of the timings that have discrepancies.
type DateOutOfAccountTimeRange struct {
	BalanceDate      time.Time
	AccountTimeRange gohtime.Range
}

// Error ensures that DateOutOfAccountTimeRange adheres to the error interface.
func (e DateOutOfAccountTimeRange) Error() string {
	return fmt.Sprintf("Balance Date is outside of Account Time Range.")
}