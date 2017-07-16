package GOHMoney

import (
	"github.com/lib/pq"
)

// NullTime is a wrapper to a pq.NullTime object to extend functionality by adding more methods
type NullTime pq.NullTime

// equal returns true if the two NullTime objects are exactly the same.
// equal even evaluates the Time fields of both NullTime objects if they are both not Valid
func (a NullTime) equal(b NullTime) bool {
	if a.Valid != b.Valid || !a.Time.Equal(b.Time) {
		return false
	}
	return true
}