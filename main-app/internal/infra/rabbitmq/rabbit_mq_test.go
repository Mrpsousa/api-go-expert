package rabbitmq

import (
	"encoding/json"
	"reflect"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/stretchr/testify/assert"
)

func Prepare() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func TestPublisher(t *testing.T) {
	exchange := "mysqlEx"
	routing_key := "product"
	product := `"id":"574f5a09-5fc7-4e8c-921b-e63c9d318c39","name":"My Product-1","price":33,"created_at":"2023-12-28T12:08:05.436225697-03:00"`

	conn, ch, err := Prepare()
	assert.Nil(t, err)

	defer conn.Close()
	defer ch.Close()

	rbClient := NewRabbit(ch)
	assert.NotNil(t, rbClient)
	assert.Equal(t, reflect.TypeOf(rbClient), reflect.TypeOf(&Rabbit{}))

	productJsonString, err := json.Marshal(product)
	assert.Nil(t, err)

	err = rbClient.Publisher(exchange, routing_key, productJsonString)
	assert.Nil(t, err)
}
