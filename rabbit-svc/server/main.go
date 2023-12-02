package main

import (
	"log"

	rb "svc/rabbitMq.com/internal/infra/rabbitmq"
)

const (
	address    = "amqp://guest:guest@localhost"
	exchange   = "mysqlEx"
	routingKey = "product"
	queueName  = "mysql"
)

func main() {
	rabbitmq, err := rb.NewRabbitMq(address)
	if err != nil {
		log.Fatal(err)
	}

	rabbitmq.Consumer(exchange, queueName, routingKey)
}
