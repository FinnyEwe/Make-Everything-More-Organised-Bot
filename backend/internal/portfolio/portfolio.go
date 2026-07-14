package portfolio

import (
	"log"

	"backend/internal/clients"
	"backend/internal/config"
)

func BuildMessage(cfg *config.Config) string {
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

func TotalValueIncrease(stock clients.Stock) (float64, float64) {
	dollar := (stock.Price - stock.CostBasis) * stock.Units
	percen := ((stock.Price - stock.CostBasis) / stock.CostBasis) * 100
	return dollar, percen
}
