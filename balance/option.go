package balance

// Option is a function that takes a pointer to a Balance returning an error.
// The idea of Option is to alter a Balance object
type Option func(*Balance) error

// Amount is an Option that will alter the amount of a Balance object.
func Amount(a int64) Option {
	return func(b *Balance) error {
		b.amount = a
		return nil
	}
}
