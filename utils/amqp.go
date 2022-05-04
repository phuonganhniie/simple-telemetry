package utils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func ConnectAmqp(user, pass, host, port string) (*amqp.Channel, func() error) {
	// connect to rabbitmq server
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	conn, err := amqp.Dial(address)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare("exchange", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, conn.Close
}
