package portfolio

import (
	"fmt"
	"strings"

	"backend/internal/clients"
)

func FormatMessage(total float64, stockList []clients.Stock, increases map[string]float64) string {
	var b strings.Builder

	totalDailyDollar := 0.0
	for _, stock := range stockList {
		p := increases[stock.Symbol]
		dailyDollar := stock.Value - stock.Value/(1+p/100)
		totalDailyDollar += dailyDollar

		b.WriteString(fmt.Sprintf(
			"**%s**  %s  `$%+.2f (%+.2f%%)`\n",
			stock.Symbol,
			formatWholeDollars(stock.Value),
			dailyDollar,
			p,
		))
	}

	totalDailyPercent := 0.0
	if startValue := total - totalDailyDollar; startValue != 0 {
		totalDailyPercent = (totalDailyDollar / startValue) * 100
	}

	b.WriteString(fmt.Sprintf(
		"**Total** %s · `$%+.2f (%+.2f%%)`\n",
		formatWholeDollars(total),
		totalDailyDollar,
		totalDailyPercent,
	))

	return b.String()
}

func formatWholeDollars(amount float64) string {
	return "$" + formatCommas(fmt.Sprintf("%.0f", amount))
}

func formatCommas(s string) string {
	parts := strings.Split(s, ".")
	intPart := parts[0]

	negative := strings.HasPrefix(intPart, "-")
	if negative {
		intPart = intPart[1:]
	}

	var b strings.Builder
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteRune(',')
		}
		b.WriteRune(c)
	}

	result := b.String()
	if negative {
		result = "-" + result
	}
	if len(parts) > 1 {
		result += "." + parts[1]
	}

	return result
}
