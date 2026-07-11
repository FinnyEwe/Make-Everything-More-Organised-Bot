package main

import (
	"log"
	"time"
	"os"
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

		// sess.AddHandler()

		if sess != nil {
			
		}
	

	now := time.Now().Format("15:04:05")	
	log.Print("hi" + now)

}