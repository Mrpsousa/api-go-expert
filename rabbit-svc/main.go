package main

import (
	"log"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

const (
	address    = "amqp://guest:guest@localhost"
	exchange   = "mysqlEx"
	routingKey = "product"
	queueName  = "mysql"
)

func main() {
	conn, err := rabbitmq.NewConn(
		address,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(exchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
}
