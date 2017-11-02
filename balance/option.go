package balance

// Option is a function that takes a pointer to a balance returning an error.
// The idea of Option is to alter a balance object
type Option func(*balance) error

// amount is an Option that will alter the amount of a balance object.
func Amount(a int) Option {
	return func(b *balance) error {
		b.amount = a
		return nil
	}
}
