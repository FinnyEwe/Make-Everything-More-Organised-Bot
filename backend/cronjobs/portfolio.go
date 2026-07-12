package cronjobs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/bwmarrin/discordgo"
	"encoding/json"
	"net/url"

	"github.com/joho/godotenv"
)

type Stock struct {
	Symbol    string
	Units     float64
	Price     float64
	CostBasis float64
	Value     float64
}

func GrabPortfolio(
// sess *discordgo.Session, message *discordgo.MessageCreate
) {
	now := time.Now().Format("15:04:05")

	total, stockList := fetchTotalPriceAndStockList()
	fmt.Print(total)
	fmt.Print(stockList)

	if now == "09:00:00" {
		// fetch macquarie
		// fetch portfolio

		// daily increase of each $% 

		// curr price of holding | done
		// total $% | done 
	}

}

func totalValueIncrease(stock Stock) (float64, float64) {
	dollar := (stock.Price - stock.CostBasis) * stock.Units
	percen := ((stock.Price - stock.CostBasis) / stock.CostBasis) * 100
	return dollar, percen
}

func fetchTotalPriceAndStockList() (float64, []Stock) {

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

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{}
	query.Set("clientId", os.Getenv("SNAPTRADE_CLIENT_ID"))
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	rawQuery := query.Encode()
	req, _ := http.NewRequest(http.MethodGet, "https://api.snaptrade.com/api/v1/accounts/"+os.Getenv("STAKE_ID")+"/positions/all?"+rawQuery, nil)

	signature := createSignature(nil, "/api/v1/accounts/"+os.Getenv("STAKE_ID")+"/positions/all", rawQuery)

	req.Header.Set("Signature", signature)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

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

func createSignature(content any, path string, query string) string {
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
	mac := hmac.New(sha256.New, []byte(os.Getenv("SNAPTRADE_CONSUMER_KEY")))
	mac.Write([]byte(strings.TrimSuffix(buf.String(), "\n")))

	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}
