package clients

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"backend/internal/config"
)

func FetchPositions(cfg *config.Config) (float64, []Stock) {
	var data struct {
		Results []struct {
			Instrument struct {
				RawSymbol string `json:"raw_symbol"`
			}
			Units     string
			Price     string
			CostBasis string `json:"cost_basis"`
		}
	}

	query := url.Values{}
	query.Set("clientId", cfg.SnapTradeClientID)
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	rawQuery := query.Encode()
	path := "/api/v1/accounts/" + cfg.StakeID + "/positions/all"
	req, _ := http.NewRequest(http.MethodGet, "https://api.snaptrade.com"+path+"?"+rawQuery, nil)

	signature := createSignature(nil, path, rawQuery, cfg.SnapTradeConsumerKey)

	req.Header.Set("Signature", signature)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&data)

	total := 0.0
	var stockList []Stock

	for _, stock := range data.Results {
		units, _ := strconv.ParseFloat(stock.Units, 64)
		price, _ := strconv.ParseFloat(stock.Price, 64)
		costBasis, _ := strconv.ParseFloat(stock.CostBasis, 64)
		value := units * price

		stockList = append(stockList, Stock{
			Symbol:    stock.Instrument.RawSymbol,
			Units:     units,
			Price:     price,
			CostBasis: costBasis,
			Value:     value,
		})
		total += value
	}

	return total, stockList
}

func createSignature(content any, path string, query string, consumerKey string) string {
	sigObject := map[string]interface{}{
		"content": content,
		"path":    path,
		"query":   query,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	sorted := map[string]interface{}{
		"content": sigObject["content"],
		"path":    sigObject["path"],
		"query":   sigObject["query"],
	}

	enc.Encode(sorted)
	mac := hmac.New(sha256.New, []byte(consumerKey))
	mac.Write([]byte(strings.TrimSuffix(buf.String(), "\n")))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
