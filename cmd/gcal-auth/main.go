package  main
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)
func main() {
	_ = godotenv.Load()
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET in .env")
	}
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob", // paste-code / desktop flow
		Scopes:       []string{calendar.CalendarScope},
	}
	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("Open this URL in your browser:")
	fmt.Println(authURL)
	fmt.Print("\nPaste the code here: ")
	var code string
	if _, err := fmt.Scanln(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("token exchange failed: %v", err)
	}
	b, _ := json.MarshalIndent(tok, "", "  ")
	fmt.Println("\nToken:")
	fmt.Println(string(b))
	fmt.Println("\nAdd this to .env:")
	fmt.Printf("GOOGLE_REFRESH_TOKEN=%s\n", tok.RefreshToken)
}