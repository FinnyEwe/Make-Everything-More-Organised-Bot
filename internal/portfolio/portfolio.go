package portfolio

import (
	"fmt"
	"log"

	"backend/internal/clients"
	"backend/internal/config"
	"backend/internal/store"
)

func BuildStockMessage(cfg *config.Config) string {
	total, stockList := clients.FetchPositions(cfg)

	var tickerSymbols []string
	for _, stock := range stockList {
		tickerSymbols = append(tickerSymbols, stock.Symbol+".AU")
	}

	increases := clients.FetchDailyIncrease(cfg, tickerSymbols)
	if increases == nil {
		log.Fatal("No daily increases found")
	}

	return FormatMessage(total, stockList, increases)
}

func BuildTotalMessage(cfg *config.Config, st *store.Store) string {
	total, _ := clients.FetchPositions(cfg)
	savings := st.GetSavings()

	grandTotal := total + savings.Amount
	// 4. Format the message
	message := fmt.Sprintf("**Total**: `$%.2f`\n**Savings**: `$%.2f`\n**Stocks**: `$%.2f`",
		grandTotal,
		savings.Amount,
		total,
	)

	return message
}
