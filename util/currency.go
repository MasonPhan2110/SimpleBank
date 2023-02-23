package util

const (
	USD = "USD"
	VND = "VND"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, VND:
		return true
	}
	return false
}
