package rabbitMQ

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	//connect to RabbitMQ
	var counts int64

	var backOff = 1 * time.Second

	var connection *amqp.Connection

	//don't continue until we have a connection
	for connection == nil {
		var err error
		connection, err = amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ. Retrying in %s second", backOff)
			time.Sleep(backOff)
			backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
			counts++
			if counts > 5 {
				log.Fatal("Failed to connect to RabbitMQ. Max retries exceeded")
				return nil, err
			}
		}
		continue
	}
	return connection, nil
}
