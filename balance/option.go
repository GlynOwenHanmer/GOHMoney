package balance

type Option func(*Balance) error

func Amount(a int64) Option {
	return func(b *Balance) error {
		b.amount = a
		return nil
	}
}