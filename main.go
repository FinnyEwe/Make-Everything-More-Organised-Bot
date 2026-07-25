package main

import (
	"backend/internal/config"
	"backend/internal/discord"
	"backend/internal/store"
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
	if err != nil {
		log.Fatal(err)
	}
	st := store.NewStore(db)

	sess, err := discordgo.New(cfg.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}

	if err := sess.Open(); err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
if err := discord.RegisterSavingsCommands(sess, st); err != nil {
		log.Fatal(err)
	}

	
	discord.GrabTotals(sess, cfg, st)
	discord.GrabPortfolio(sess, cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
