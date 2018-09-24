package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel
var queue = flag.Int("queue", 1, "-queue=1")

func main() {
	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.WithError(err).Fatal("connection")
	}
	defer conn.Close()
	ch, err = conn.Channel()
	if err != nil {
		logrus.WithError(err).Fatal("channel")
	}
	defer ch.Close()

	qName := fmt.Sprintf("queue-%d", *queue)
	messages, err := ch.Consume(
		qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	fmt.Printf("listening on queue=%s\n", qName)
	for msg := range messages {
		logrus.Infof("GOT: %s", string(msg.Body))
	}
}
