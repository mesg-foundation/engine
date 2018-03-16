package conversion

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Currency type definition
type Currency uint

// Unit association with number of zero
const (
	BASE Currency = 1 // base value
	MESG Currency = 1000000000
)

var stringToCurrency = map[string]Currency{
	"":     MESG, // default currency
	"MESG": MESG,
	"BASE": BASE,
}

var regex *regexp.Regexp

// Amount struct
type Amount struct {
	Value uint
}

func stringFromCurrency(currency Currency) (value string) {
	for key := range stringToCurrency {
		if key != "" && stringToCurrency[key] == currency {
			return key
		}
	}
	return ""
}

func init() {
	// Build regex
	currencies := make([]string, 0, len(stringToCurrency))
	for key := range stringToCurrency {
		currencies = append(currencies, key)
	}
	currenciesString := strings.Join(currencies, "|")
	regex = regexp.MustCompile("^([0-9]*[.]?[0-9]+)( (" + currenciesString + "))?$")
}

// extractFromString extracts a number and a currency from a string
func extractFromString(value string) (number string, currency string, error error) {
	result := regex.FindStringSubmatch(value)
	if result == nil || len(result) != 4 {
		error = errors.New("Currency string is not valid. Should be like '42.42 MESG'")
		return
	}
	number, currency = result[1], result[3]
	return
}

// toBase converts a number with a currency to the base value
func toBase(numberString string, currencyString string) (base float64, err error) {
	number, err := strconv.ParseFloat(numberString, 64)
	currency := stringToCurrency[currencyString]
	base = number * float64(currency)
	return
}

// FromString converts a string containing a value and an currency unit into the base value
func (amount *Amount) FromString(value string) (err error) {
	numberString, currencyString, err := extractFromString(value)
	if err != nil {
		return
	}
	base, _ := toBase(numberString, currencyString)
	amount.Value = uint(base)
	if float64(amount.Value) != base {
		err = errors.New("The number lost some precision. Check your value and unit")
	}
	return
}

// Convert converts an Amount struct to a specific Currency
func (amount *Amount) Convert(currency Currency) (value float64) {
	value = float64(amount.Value) / float64(currency)
	return
}

// Format returns a nice human readable amount using a given currency
func (amount *Amount) Format(currency Currency) (desc string) {
	currencyStr := stringFromCurrency(currency)
	desc = strconv.FormatFloat(amount.Convert(currency), 'f', -1, 64) + " " + currencyStr
	return
}

// String displays the amount using the MESG currency
func (amount *Amount) String() (desc string) {
	desc = amount.Format(MESG)
	return
}
