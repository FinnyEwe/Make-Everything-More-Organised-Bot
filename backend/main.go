package main

import (
	"log"
	"os"
	"os/signal"

	"backend/internal/config"
	"backend/internal/discord"

	discordgo "github.com/bwmarrin/discordgo"
)

func main() {
	cfg := config.Load()
	sess, err := discordgo.New(cfg.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}

	if err := sess.Open(); err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	discord.GrabPortfolio(sess, cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
