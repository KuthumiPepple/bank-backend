package util

const (
	CAD = "CAD"
	EUR = "EUR"
	GBP = "GBP"
	NGN = "NGN"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CAD, EUR, GBP, NGN, USD:
		return true
	}
	return false
}
