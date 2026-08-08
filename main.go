package main

import (
	"backend/internal/config"
	"backend/internal/discord"
	"backend/internal/reminders"
	"backend/internal/store"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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

	// Setup RabbitMQ for reminders
	ch, err := reminders.CreateQueueConnection()
	if err != nil {
		log.Fatal(err)
	}
	reminders.CreateExchangeQueue(ch)

	// Start reminder consumers
	go reminders.ConsumeWork(ch, st, sess, cfg.DiscordChannelID)

	// Setup cron jobs for 8:30 AM daily (Sydney time)
	c := cron.New(cron.WithLocation(mustLoadLocation("Australia/Sydney")))
	
	// Portfolio updates at 8:30 AM
	c.AddFunc("30 8 * * *", func() {
		log.Println("Running scheduled portfolio updates...")
		discord.GrabTotals(sess, cfg, st)
		discord.GrabPortfolio(sess, cfg)
	})

	// Reminder notifications at 8:30 AM
	c.AddFunc("30 8 * * *", func() {
		log.Println("Running scheduled reminder check...")
		reminders.QueueNotifications(st, ch)
	})

	// Cleanup completed reminders at 9:00 AM (after processing)
	c.AddFunc("0 9 * * *", func() {
		log.Println("Cleaning up completed reminders...")
		if err := st.DeleteCompletedReminders(); err != nil {
			log.Printf("Failed to delete completed reminders: %v", err)
		} else {
			log.Println("Successfully cleaned up completed reminders")
		}
	})

	c.Start()
	defer c.Stop()

	log.Println("Bot started. Scheduled jobs will run at 8:30 AM Sydney time.")
	log.Println("Press Ctrl+C to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Fatalf("Failed to load timezone %s: %v", name, err)
	}
	return loc
}
