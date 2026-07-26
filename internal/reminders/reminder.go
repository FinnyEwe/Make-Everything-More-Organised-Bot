package reminders

import (
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
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

func ConsumeWork(ch *amqp.Channel){
	discChan, _ := ch.Consume("discord",  
	"",     // Auto-ack?
    false,     // Exclusive?
    false,     // No-local (not used by RabbitMQ)
    false, 
	false,    // No-wait?
    nil,       // Arguments)
	)

	gcalChan, _ := ch.Consume("gcal",  
	"",     // Auto-ack?
    false,     // Exclusive?
    false,     // No-local (not used by RabbitMQ)
    false, 
	false,    // No-wait?
    nil,       // Arguments)
	)

	go SendDiscord(discChan)
	go SendCal(gcalChan)	
}


func SendDiscord(messages  <-chan amqp.Delivery){
	for d := range messages { // 'msgs' is the <-chan amqp.Delivery from ch.Consume
		//get reminder by id
		//send to discord
		// tick off in db
	d.Ack(false)
}
}

func SendCal(messages  <-chan amqp.Delivery){
	for d := range messages { // 'msgs' is the <-chan amqp.Delivery from ch.Consume
	//get reminder by id
	//send to gcal
	// tick off in db
	d.Ack(false)
}
}

func CreateQueueConnection() (*amqp.Channel, error) {
	godotenv.Load()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:%s@rabbitmq:5672/", os.Getenv("RABBITMQ_PASS")))
	failOnError(err, "failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
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