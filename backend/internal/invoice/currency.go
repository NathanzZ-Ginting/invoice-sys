package invoice

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// GetCurrencySymbol returns the symbol for a given currency code
// For PDF compatibility, uses ASCII-safe symbols or currency codes
func GetCurrencySymbol(currencyCode string) string {
	symbols := map[string]string{
		"USD": "$",
		"IDR": "Rp",
		"EUR": "EUR ", // Use code for PDF compatibility
		"GBP": "GBP ", // Use code for PDF compatibility
		"JPY": "JPY ", // Use code for PDF compatibility
		"CNY": "CNY ", // Use code for PDF compatibility
		"SGD": "S$",
		"MYR": "RM",
	}
	
	if symbol, ok := symbols[currencyCode]; ok {
		return symbol
	}
	return currencyCode + " "
}

// FormatCurrency formats amount based on currency type
func FormatCurrency(amount float64, currencyCode string) string {
	// Currencies without decimals
	noDecimalCurrencies := map[string]bool{
		"IDR": true,
		"JPY": true,
		"CNY": true,
	}
	
	p := message.NewPrinter(language.English)
	
	if noDecimalCurrencies[currencyCode] {
		// Round and format without decimals
		return p.Sprintf("%.0f", amount)
	}
	
	// Format with 2 decimals
	return p.Sprintf("%.2f", amount)
}

// FormatCurrencyWithSymbol formats amount with currency symbol
func FormatCurrencyWithSymbol(amount float64, currencyCode string) string {
	symbol := GetCurrencySymbol(currencyCode)
	formatted := FormatCurrency(amount, currencyCode)
	return symbol + formatted
}
