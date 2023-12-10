package rabbitmq

type RabbitCHInterface interface {
	Publisher(eX, routKey string, productJsonString []byte) error
}
