package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// The first value from the slice `os.Args` is the path to the program.
	args := os.Args[1:]
	fmt.Println("Hello, World!")

	if len(args) < 2 {
		fmt.Println(`Need at least two arguments:
1st positional argument: url connection to RabbitMQ
2nd positional argument: rabbitmq queue name to send message to`)
		panic("Expect at least two arguments")
	}

	rabbitMqQueue := args[1]
	rabbitMqConnUrl := args[0]

	conn, err := amqp.Dial(rabbitMqConnUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// The queue will only be created if it doesn't exist already.
	q, err := ch.QueueDeclare(
		rabbitMqQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgBody := "Finished."

	// Simulate working.
	time.Sleep(5 * time.Second)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msgBody),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", msgBody)
}
