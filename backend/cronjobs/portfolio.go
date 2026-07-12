package cronjobs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
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

func GrabPortfolio(
// sess *discordgo.Session, message *discordgo.MessageCreate
) {
	now := time.Now().Format("15:04:05")
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{}
	query.Set("clientId", os.Getenv("SNAPTRADE_CLIENT_ID"))
	query.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	rawQuery := query.Encode()
	req, _ := http.NewRequest(http.MethodGet, "https://api.snaptrade.com/api/v1/accounts?"+rawQuery, nil)
	sigObject := map[string]interface{}{
		"content": nil,
		"path":    "/api/v1/accounts",
		"query":   rawQuery,
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

	req.Header.Set("Signature", signature)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

	fmt.Print(resp)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("status:", resp.Status)
	fmt.Println("body:", string(body))








	if now == "09:00:00" {
		// fetch macquarie
		// fetch portfolio

		// daily increase of each
		// total

	}

}
