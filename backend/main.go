package main

import (
	"log"
	"os"
	"os/signal"

	"backend/cronjobs"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)
func main(){
	err := godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	sess, err := discordgo.New(token)
	if err != nil {
		log.Fatal(err)
	}

	if err := sess.Open(); err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	cronjobs.GrabPortfolio(sess)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

}