package rabbitmq

import (
	"log"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

// address="amqp://guest:guest@localhost"
// exchange="events"
// routing_key="product"
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

func (r *RabbitMq) Publisher(exchange, routing_key, msq string) {
	defer r.RMqConn.Close()

	publisher, err := rabbitmq.NewPublisher(
		r.RMqConn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	err = publisher.Publish(
		[]byte(msq),
		[]string{routing_key},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchange),
	)
	if err != nil {
		log.Println(err)
	}
}
