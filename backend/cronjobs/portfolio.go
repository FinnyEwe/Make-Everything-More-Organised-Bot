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

	"encoding/json"
	"net/url"

	"github.com/bwmarrin/discordgo"

	"github.com/joho/godotenv"
)

type Stock struct {
	Symbol    string
	Units     float64
	Price     float64
	CostBasis float64
	Value     float64
}

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

func GrabPortfolio(sess *discordgo.Session) {
	// now := time.Now().Format("15:04:05")

	total, stockList := fetchTotalPriceAndStockList()

	var tickerSymbols []string
	for _, stock := range stockList {
		tickerSymbols = append(tickerSymbols, stock.Symbol+".AU")
	}

	increases := dailyIncrease(tickerSymbols)
	if increases == nil {
		log.Fatal("No daily increases found")
	}

	message := formatPortfolioMessage(total, stockList, increases)
	fmt.Println(message)

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		log.Println("DISCORD_CHANNEL_ID not set, printing message:")
		fmt.Println(message)
		return
	}
	_, err := sess.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Fatal(err)
	}
}

func formatPortfolioMessage(total float64, stockList []Stock, increases map[string]float64) string {
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

func dailyIncrease(symbolList []string) map[string]float64 {
	if len(symbolList) == 0 {
		return nil
	}

	rawQuery := url.Values{}
	rawQuery.Set("api_token", os.Getenv("EODHD_API_KEY"))
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
