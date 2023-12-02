package rabbitmq

import (
	"log"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

type RabbitMq struct {
	RMqConn *rabbitmq.Conn
}

func NewRabbitMq(address string) (*RabbitMq, error) {
	conn, err := rabbitmq.NewConn(
		address,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMq{RMqConn: conn}, nil
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
