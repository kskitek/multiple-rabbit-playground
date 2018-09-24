package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel
var noOfQueues = flag.Int("queues", 1, "-queues=1")
var noOfExchanges = flag.Int("exchanges", 2, "-exchanges=2")

func main() {
	flag.Parse()
	fmt.Printf("exchanges=%d, queues=%d\n", *noOfExchanges, *noOfQueues)

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

	// TODO
	//err = ch.Qos(
	//	1,
	//	0,
	//	false,
	//)

	for ex := 1; ex <= *noOfExchanges; ex++ {
		for q := 1; q <= *noOfQueues; q++ {
			fmt.Println("Declaring", "ex", ex, "q", q)
			declareQueue(ex, q)
		}
	}

	fmt.Println("PUBLISHING")
	for {
		var exNo int
		fmt.Print("Ex number: ")
		fmt.Scanf("%d", &exNo)
		publish(exNo)
		fmt.Println("")
	}

}

func publish(exNo int) {
	exName := fmt.Sprintf("ex-%d", exNo)
	err := ch.Publish(
		exName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("MESSAGE"),
		})
	if err != nil {
		logrus.WithError(err).Fatal("publish")
	}
}

func declareQueue(exNo, qNo int) {
	exName := fmt.Sprintf("ex-%d", exNo)
	err := ch.ExchangeDeclare(
		exName,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.WithError(err).Fatal("exchange")
	}

	qName := fmt.Sprintf("queue-%d", qNo)
	_, err = ch.QueueDeclare(
		qName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.WithError(err).Fatal("queue")
	}

	err = ch.QueueBind(
		qName,
		"",
		exName,
		false,
		nil,
	)
	if err != nil {
		logrus.WithError(err).Fatal("bind")
	}
}
