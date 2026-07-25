package discord

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/portfolio"
	"backend/internal/store"

	"github.com/bwmarrin/discordgo"
)

func GrabPortfolio(sess *discordgo.Session, cfg *config.Config) {
	message := portfolio.BuildStockMessage(cfg)

	fmt.Println(message)

	if cfg.DiscordChannelID == "" {
		log.Println("DISCORD_CHANNEL_ID not set, printing message:")
		fmt.Println(message)
		return
	}
	_, err := sess.ChannelMessageSend(cfg.DiscordChannelID, message)
	if err != nil {
		log.Fatal(err)
	}
}

func GrabTotals(sess *discordgo.Session, cfg *config.Config, st *store.Store) {
	message := portfolio.BuildTotalMessage(cfg, st)
	if cfg.DiscordChannelID == "" {
		log.Println("DISCORD_CHANNEL_ID not set, printing message:")
		fmt.Println(message)
		return
	}
	_, err := sess.ChannelMessageSend(cfg.DiscordChannelID, message)
	if err != nil {
		log.Fatal(err)
	}
}
