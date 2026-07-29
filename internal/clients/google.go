package clients

import (
	"context"
	"os"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)
func NewCalendarService(ctx context.Context) (*calendar.Service, error) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarScope},
	}
	token := &oauth2.Token{
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
	}
	httpClient := config.Client(ctx, token)
	return calendar.NewService(ctx, option.WithHTTPClient(httpClient))
}