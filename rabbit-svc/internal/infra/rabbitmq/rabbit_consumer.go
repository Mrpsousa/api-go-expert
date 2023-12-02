package rabbitmq

import (
	"log"

	"svc/rabbitMq.com/internal/infra/database"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

type RabbitMq struct {
	RMqConn   *rabbitmq.Conn
	ProductDB database.ProductInterface
}

func NewRabbitMq(address string, db database.ProductInterface) (*RabbitMq, error) {
	conn, err := rabbitmq.NewConn(
		address,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMq{RMqConn: conn,
		ProductDB: db}, nil
}

func (r *RabbitMq) Consumer(exchange, queueName, routingKey string) {
	defer r.RMqConn.Close()

	consumer, err := rabbitmq.NewConsumer(
		r.RMqConn,
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
