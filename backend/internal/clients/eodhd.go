package clients

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"backend/internal/config"
)

type DailyIncrease struct {
	Change        float64
	ChangePercent float64 `json:"change_p"`
	Close         float64
	Code          string
	GmtOffset     float64
	High          float64
	Low           float64
	Open          float64
	PreviousClose float64
	Timestamp     float64
	Volume        float64
}

func FetchDailyIncrease(cfg *config.Config, symbolList []string) map[string]float64 {
	if len(symbolList) == 0 {
		return nil
	}

	rawQuery := url.Values{}
	rawQuery.Set("api_token", cfg.EODHDAPIKey)
	if len(symbolList) > 1 {
		rawQuery.Set("s", strings.Join(symbolList[1:], ","))
	}
	rawQuery.Set("fmt", "json")

	req, _ := http.NewRequest(http.MethodGet, "https://eodhd.com/api/real-time/"+symbolList[0]+"?"+rawQuery.Encode(), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var data []DailyIncrease
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	increases := make(map[string]float64)
	for _, stock := range data {
		increases[strings.Split(stock.Code, ".")[0]] = stock.ChangePercent
	}

	return increases
}
