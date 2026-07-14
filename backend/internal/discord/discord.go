package discord

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/portfolio"

	"github.com/bwmarrin/discordgo"
)

func GrabPortfolio(sess *discordgo.Session, cfg *config.Config) {
	message := portfolio.BuildMessage(cfg)

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
