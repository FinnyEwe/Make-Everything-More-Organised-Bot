package cronjobs

import (
	"time"
	"github.com/bwmarrin/discordgo"
) 

func grabPortfolio(sess *discordgo.Session, message *discordgo.MessageCreate){
	now := time.Now().Format("15:04:05")

	if now == "09:00:00" {
		// fetch macquarie
		// fetch portfolio

		// daily increase of each
		// total 
	}
	
}