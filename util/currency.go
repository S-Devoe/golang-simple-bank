package util

const (
	USD = "USD"
	NGN = "NGN"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, NGN, EUR, CAD:
		return true
	}
	return false
}
