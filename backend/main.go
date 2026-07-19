package main

import (
	"backend/internal/config"
	"backend/internal/discord"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	print(db)

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
