package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Ch *amqp.Channel
}

func NewRabbit(ch *amqp.Channel) *Rabbit {
	return &Rabbit{Ch: ch}
}

func (rb *Rabbit) Publisher(eX, routKey string, productJsonString []byte) error {

	err := rb.Ch.PublishWithContext(context.TODO(),
		eX,      // exchange
		routKey, // routing key
		false,   // ch
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        productJsonString,
		})

	return err
}
