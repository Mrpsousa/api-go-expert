package rabbitmq

import (
	"encoding/json"
	"log"

	"svc/rabbitMq.com/internal/entity"
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
	product := &entity.Product{}

	consumer, err := rabbitmq.NewConsumer(
		r.RMqConn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			// Verificar se o corpo está presente
			if len(d.Body) == 0 {
				log.Println("Received empty message, skipping.")
				return rabbitmq.Ack
			}

			// O corpo está presente, prosseguir com a deserialização
			err := json.Unmarshal(d.Body, product)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				return rabbitmq.NackDiscard
			}

			// Aqui você pode processar o produto, armazená-lo no banco de dados, etc.
			r.ProductDB.Create(product)
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
