package reminders

import (
	"backend/internal/clients"
	"backend/internal/model"
	"backend/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/api/calendar/v3"
)

//polling is done
//schedule it will be outside
func QueueNotifications(s *store.Store, ch *amqp.Channel){
	reminders, err := s.PollReminders()
	if err != nil {
		log.Fatal(err)
	}

	targets := []string{"discord", "gcal"}

	for _, reminder := range reminders {
		for _, target := range targets {

			body, _ := json.Marshal(reminder.ID)

			ch.Publish(
				"notifications",
				target,
				false,
				false,
				amqp.Publishing{
					Body: body,
				},
			)
		}
	}

}

func ConsumeWork(ch *amqp.Channel, s *store.Store, sess *discordgo.Session, channelID string) {
	discChan, err := ch.Consume(
		"discord",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("discord consume failed: %v", err)
		return
	}

	gcalChan, err := ch.Consume(
		"gcal",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("gcal consume failed: %v", err)
		return
	}

	go SendDiscord(discChan, s, sess, channelID)
	go SendCal(gcalChan, s)
}

func SendDiscord(messages <-chan amqp.Delivery, s *store.Store, sess *discordgo.Session, channelID string) {
	loc, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		log.Printf("load location failed: %v", err)
		return
	}

	for d := range messages {
		var id uint
		if err := json.Unmarshal(d.Body, &id); err != nil {
			log.Printf("discord: bad message body: %v", err)
			d.Nack(false, false)
			continue
		}

		var reminder model.Reminder
		if err := s.Db.Where("id = ?", id).First(&reminder).Error; err != nil {
			log.Printf("discord: reminder %d not found: %v", id, err)
			d.Nack(false, false)
			continue
		}

		if reminder.Discord {
			d.Ack(false)
			continue
		}

		msg := fmt.Sprintf("⏰ **Reminder** — %s\n%s", reminder.Date, reminder.Description)
		if _, err := sess.ChannelMessageSend(channelID, msg); err != nil {
			log.Printf("discord: send failed for reminder %d: %v", id, err)
			d.Nack(false, true)
			continue
		}

		if time.Now().In(loc).Format("02-01-2006") == reminder.Date {
			s.Db.Model(&model.Reminder{}).Where("id = ?", id).Update("discord", true)
		}

		d.Ack(false)
	}
}

func SendCal(messages  <-chan amqp.Delivery, s *store.Store){

	svc, err := clients.NewCalendarService(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		
	for d := range messages {
		var id uint
		json.Unmarshal(d.Body, &id)
		//send to gcal

		var reminder model.Reminder


		s.Db.Where("id = ?", id).First(&reminder)

		loc, err := time.LoadLocation("Australia/Sydney")
		if err != nil {
			log.Fatal(err)
		}

		day, err := time.ParseInLocation("02-01-2006", reminder.Date, loc)
		if err != nil {
			log.Printf("bad date %q: %v", reminder.Date, err)
			d.Nack(false, false)
			continue
		}
		
		start := time.Date(day.Year(), day.Month(), day.Day(), 8, 30, 0, 0, loc)

		event := &calendar.Event{
			Summary: reminder.Description,
			Start: &calendar.EventDateTime{
				DateTime: start.Format(time.RFC3339), 
				TimeZone: "Australia/Sydney",
			},
			End: &calendar.EventDateTime{
				DateTime: start.Format(time.RFC3339),
				TimeZone: "Australia/Sydney",
			},
		}
		created, err := svc.Events.Insert("primary", event).Do()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(created.Id, created.HtmlLink)


		// tick off in db
		
		if time.Now().In(loc).Equal(day) {
			s.Db.Where("id = ?", id).Update("gcal", true)
		}
		


	d.Ack(false)
}
}

func CreateQueueConnection() (*amqp.Channel, error) {
	godotenv.Load()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest1:%s@localhost:5672/", os.Getenv("RABBITMQ_PASS")))
	failOnError(err, "failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch, err

}

func CreateExchangeQueue(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"notifications", 
		"direct",        
		true,          
		false,          
		false,           
		false,           
		nil,             
	)
	failOnError(err, "Failed to open a exchange")

	gcalQueue, _ := ch.QueueDeclare(
		"gcal",
		true,
		false,
		false,
		false,
		nil,
	)

	discordQueue, _ := ch.QueueDeclare(
		"discord",
		true,
		false,
		false,
		false,
		nil,
	)

	ch.QueueBind(
		gcalQueue.Name,
		"gcal",
		"notifications",
		false,
		nil,
	)

	ch.QueueBind(
		discordQueue.Name,
		"discord",        // routing key
		"notifications",  // exchange
		false,
		nil,
	)
	
}



func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }