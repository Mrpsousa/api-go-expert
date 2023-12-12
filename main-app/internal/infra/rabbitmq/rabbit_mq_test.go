package rabbitmq

import (
	"reflect"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
func TestCreateUser(t *testing.T) {
	conn, ch, err := Prepare()
	assert.Nil(t, err)

	defer conn.Close()
	defer ch.Close()

	rbClient := NewRabbit(ch)
	assert.NotNil(t, rbClient)
	assert.Equal(t, reflect.TypeOf(rbClient), reflect.TypeOf(&Rabbit{}))
	var product []byte

	err = rbClient.Publisher(mock.Anything, mock.Anything, product)
	assert.Nil(t, err)

}
