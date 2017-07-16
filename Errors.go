package GOHMoney

import "bytes"

// AccountFieldError holds zero or more descriptions of things that are wrong with potential new Account items.
type AccountFieldError []string

// Error ensures that AccountFieldError adheres to the error interface.
func (e AccountFieldError) Error() string {
	var errorString bytes.Buffer
	errorString.WriteString("AccountFieldError: ")
	for i, field := range e {
		errorString.WriteString(field)
		if i < len(e) - 1 {
			errorString.WriteByte(' ')
		}
	}
	return string(errorString.String())
}

// Various error strings describing possible errors with potential new Account items.
const (
	EmptyNameError                   = "Empty name."
	ZeroDateOpenedError              = "No opened date given."
	ZeroValidDateClosedError         = "Closed date marked as valid but not set."
)


