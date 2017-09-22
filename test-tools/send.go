package main

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	url := os.Getenv("AMQP_URL")

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("connection.open: %s", err)
	}

	// This waits for a server acknowledgment which means the sockets will have
	// flushed all outbound publishings prior to returning.  It's important to
	// block on Close to not lose any publishings.
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}

	log.Printf("%s", time.Now())

	for i := 0; i < 100000; i++ {
		msg := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "application/json",
			Body:         []byte("{\"hello\": \"world\"}"),
		}

		// This is not a mandatory delivery, so it will be dropped if there are no
		// queues bound to the logs exchange.
		err = c.Publish("fishtrip", "order.state.paid", false, false, msg)
		if err != nil {
			// Since publish is asynchronous this can happen if the network connection
			// is reset or if the server has run out of resources.
			log.Fatalf("basic.publish: %v", err)
		}
	}
	log.Printf("%s", time.Now())
}
