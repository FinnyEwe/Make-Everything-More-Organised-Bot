package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken         string
	DiscordChannelID     string
	SnapTradeClientID    string
	SnapTradeConsumerKey string
	StakeID              string
	EODHDAPIKey          string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	return &Config{
		DiscordToken:         os.Getenv("DISCORD_TOKEN"),
		DiscordChannelID:     os.Getenv("DISCORD_CHANNEL_ID"),
		SnapTradeClientID:    os.Getenv("SNAPTRADE_CLIENT_ID"),
		SnapTradeConsumerKey: os.Getenv("SNAPTRADE_CONSUMER_KEY"),
		StakeID:              os.Getenv("STAKE_ID"),
		EODHDAPIKey:          os.Getenv("EODHD_API_KEY"),
	}
}
