package rabbitmq

import (
	"context"

	lg "github.com/mrpsousa/api/pkg/log"
	er "github.com/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Ch    *amqp.Channel
	Queue *amqp.Queue
}

// https://pkg.go.dev/github.com/rabbitmq/amqp091-go@v1.7.0#readme-documentation

func NewRabbit(ch *amqp.Channel) *Rabbit {
	return &Rabbit{Ch: ch}
}

func PrepareAmqp() (*amqp.Connection, *amqp.Channel, error) {
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

func (rb *Rabbit) declareExchange(exchangeName string) error {

	lg.LogMessage("DeclareExchange", "Creating Exchange", false)

	err := rb.Ch.ExchangeDeclare(
		exchangeName, // Nome da exchange
		"direct",     // Exchange type(direct)
		false,        // Durable
		false,        // Auto-delete
		false,        // Internal
		false,        // No-wait
		nil,          // Argumentos adicionais
	)
	if err != nil {
		return er.Wrap(err, "Failed to declare an exchange")
	}
	return nil
}

func (rb *Rabbit) declareQueue(queueName string) error {

	lg.LogMessage("DeclareQueue", "Creating Queue", false)

	queue, err := rb.Ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return er.Wrap(err, "Failed to declare a queue")
	}

	rb.Queue = &queue

	return nil
}

func (rb *Rabbit) bindQueueToExchange(queueName, exchangeName, routingKey string) error {

	lg.LogMessage("BindQueueToExchange", "Creating bind", false)

	err := rb.Ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // No-wait
		nil,   // Argumentos adicionais
	)
	if err != nil {
		return er.Wrap(err, "Failed to bind the queue")
	}

	return nil
}

func (rb *Rabbit) Publisher(exchangeName, routKey string, productJsonString []byte) error {
	queueName := "testQueueName"
	// Verifica se a exchange já existe
	exchangeExists, err := rb.exchangeExists(exchangeName)
	if err != nil {
		return err
	}

	// if rb.Ch == nil {
	// 	fmt.Print("algo")
	// }
	// Se a exchange não existir, declare-a
	if !exchangeExists {
		lg.LogMessage("Publisher", "Exchange does not exist", false)
		rb.declareExchange(exchangeName)
		if err != nil {
			return er.Wrap(err, "Failed to declare exchange")
		}
		rb.declareQueue(queueName)
		if err != nil {
			return er.Wrap(err, "Failed to declare queue")
		}
		rb.bindQueueToExchange(queueName, exchangeName, routKey)
		if err != nil {
			return er.Wrap(err, "Failed to bind queue")
		}
	}

	// Publica a mensagem na exchange
	err = rb.Ch.PublishWithContext(context.TODO(),
		exchangeName, // exchange
		routKey,      // routing key
		false,        // ch
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        productJsonString,
		})

	return err
}

func (rb *Rabbit) exchangeExists(exchangeName string) (bool, error) {
	// Tenta declarar a exchange sem causar erro
	err := rb.Ch.ExchangeDeclarePassive(
		exchangeName, // Nome da exchange
		"direct",     // Tipo de exchange (direta)
		false,        // Durable
		false,        // Auto-delete
		false,        // Internal
		false,        // No-wait
		nil,          // Argumentos adicionais
	)

	// Se a exchange já existir, não é um erro
	if err == nil {
		return true, nil
	}

	// Se o erro for relacionado a não existência da exchange, retorna false
	if amqpErr, ok := err.(*amqp.Error); ok && amqpErr.Code == 404 {
		return false, nil
	}

	// Outro erro ocorreu
	return false, err
}

// // func (rb *Rabbit) Publisher(exchangeName, routKey string, productJsonString []byte) error {
// // 	err := rb.Ch.PublishWithContext(context.TODO(),
// // 	exchangeName, // exchange
// // 	routKey,      // routing key
// // 	false,        // ch
// // 	false,        // immediate
// // 	amqp.Publishing{
// // 		ContentType: "text/plain",
// // 		Body:        productJsonString,
// // 	})
// // 	return err
// // }

// //	func (rb *Rabbit) DeclareExchanges(exchangeName string) error {
// //		// Declare a exchange de destino
// //		err := rb.Ch.ExchangeDeclare(
// //			exchangeName, // Nome da exchange de destino
// //			"direct",     // Tipo de exchange (direct)
// //			false,        // Durable
// //			false,        // Auto-delete
// //			false,        // Internal
// //			false,        // No-wait
// //			nil,          // Argumentos adicionais
// //		)
// //		return err
// //	}
// func main() {
// 	queueName := "job_queue"
// 	addr := "amqp://guest:guest@localhost:5672/"
// 	queue := New(queueName, addr)
// 	message := []byte("message")

// 	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*20))
// 	defer cancel()
// loop:
// 	for {
// 		select {
// 		// Attempt to push a message every 2 seconds
// 		case <-time.After(time.Second * 2):
// 			if err := queue.Push(message); err != nil {
// 				fmt.Printf("Push failed: %s\n", err)
// 			} else {
// 				fmt.Println("Push succeeded!")
// 			}
// 		case <-ctx.Done():
// 			queue.Close()
// 			break loop
// 		}
// 	}
// }

// type Client struct {
// 	queueName       string
// 	logger          *log.Logger
// 	connection      *amqp.Connection
// 	channel         *amqp.Channel
// 	done            chan bool
// 	notifyConnClose chan *amqp.Error
// 	notifyChanClose chan *amqp.Error
// 	notifyConfirm   chan amqp.Confirmation
// 	isReady         bool
// }

// const (
// 	reconnectDelay = 5 * time.Second

// 	reInitDelay = 2 * time.Second

// 	resendDelay = 5 * time.Second
// )

// var (
// 	errNotConnected  = errors.New("not connected to a server")
// 	errAlreadyClosed = errors.New("already closed: not connected to the server")
// 	errShutdown      = errors.New("client is shutting down")
// )

// // New creates a new consumer state instance, and automatically
// // attempts to connect to the server.
// func New(queueName, addr string) *Client {
// 	client := Client{
// 		logger:    log.New(os.Stdout, "", log.LstdFlags),
// 		queueName: queueName,
// 		done:      make(chan bool),
// 	}
// 	go client.handleReconnect(addr)
// 	return &client
// }

// // handleReconnect will wait for a connection error on
// // notifyConnClose, and then continuously attempt to reconnect.
// func (client *Client) handleReconnect(addr string) {
// 	for {
// 		client.isReady = false
// 		client.logger.Println("Attempting to connect")

// 		conn, err := client.connect(addr)

// 		if err != nil {
// 			client.logger.Println("Failed to connect. Retrying...")

// 			select {
// 			case <-client.done:
// 				return
// 			case <-time.After(reconnectDelay):
// 			}
// 			continue
// 		}

// 		if done := client.handleReInit(conn); done {
// 			break
// 		}
// 	}
// }

// // connect will create a new AMQP connection
// func (client *Client) connect(addr string) (*amqp.Connection, error) {
// 	conn, err := amqp.Dial(addr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	client.changeConnection(conn)
// 	client.logger.Println("Connected!")
// 	return conn, nil
// }

// // handleReconnect will wait for a channel error
// // and then continuously attempt to re-initialize both channels
// func (client *Client) handleReInit(conn *amqp.Connection) bool {
// 	for {
// 		client.isReady = false

// 		err := client.init(conn)

// 		if err != nil {
// 			client.logger.Println("Failed to initialize channel. Retrying...")

// 			select {
// 			case <-client.done:
// 				return true
// 			case <-time.After(reInitDelay):
// 			}
// 			continue
// 		}

// 		select {
// 		case <-client.done:
// 			return true
// 		case <-client.notifyConnClose:
// 			client.logger.Println("Connection closed. Reconnecting...")
// 			return false
// 		case <-client.notifyChanClose:
// 			client.logger.Println("Channel closed. Re-running init...")
// 		}
// 	}
// }

// // init will initialize channel & declare queue
// func (client *Client) init(conn *amqp.Connection) error {
// 	ch, err := conn.Channel()

// 	if err != nil {
// 		return err
// 	}

// 	err = ch.Confirm(false)

// 	if err != nil {
// 		return err
// 	}
// 	_, err = ch.QueueDeclare(
// 		client.queueName,
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	client.changeChannel(ch)
// 	client.isReady = true
// 	client.logger.Println("Setup!")

// 	return nil
// }

// // changeConnection takes a new connection to the queue,
// // and updates the close listener to reflect this.
// func (client *Client) changeConnection(connection *amqp.Connection) {
// 	client.connection = connection
// 	client.notifyConnClose = make(chan *amqp.Error, 1)
// 	client.connection.NotifyClose(client.notifyConnClose)
// }

// // changeChannel takes a new channel to the queue,
// // and updates the channel listeners to reflect this.
// func (client *Client) changeChannel(channel *amqp.Channel) {
// 	client.channel = channel
// 	client.notifyChanClose = make(chan *amqp.Error, 1)
// 	client.notifyConfirm = make(chan amqp.Confirmation, 1)
// 	client.channel.NotifyClose(client.notifyChanClose)
// 	client.channel.NotifyPublish(client.notifyConfirm)
// }

// // Push will push data onto the queue, and wait for a confirm.
// // This will block until the server sends a confirm. Errors are
// // only returned if the push action itself fails, see UnsafePush.
// func (client *Client) Push(data []byte) error {
// 	if !client.isReady {
// 		return errors.New("failed to push: not connected")
// 	}
// 	for {
// 		err := client.UnsafePush(data)
// 		if err != nil {
// 			client.logger.Println("Push failed. Retrying...")
// 			select {
// 			case <-client.done:
// 				return errShutdown
// 			case <-time.After(resendDelay):
// 			}
// 			continue
// 		}
// 		confirm := <-client.notifyConfirm
// 		if confirm.Ack {
// 			client.logger.Printf("Push confirmed [%d]!", confirm.DeliveryTag)
// 			return nil
// 		}
// 	}
// }

// // UnsafePush will push to the queue without checking for
// // confirmation. It returns an error if it fails to connect.
// // No guarantees are provided for whether the server will
// // receive the message.
// func (client *Client) UnsafePush(data []byte) error {
// 	if !client.isReady {
// 		return errNotConnected
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	return client.channel.PublishWithContext(
// 		ctx,
// 		"",
// 		client.queueName,
// 		false,
// 		false,
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        data,
// 		},
// 	)
// }

// // Consume will continuously put queue items on the channel.
// // It is required to call delivery.Ack when it has been
// // successfully processed, or delivery.Nack when it fails.
// // Ignoring this will cause data to build up on the server.
// func (client *Client) Consume() (<-chan amqp.Delivery, error) {
// 	if !client.isReady {
// 		return nil, errNotConnected
// 	}

// 	if err := client.channel.Qos(
// 		1,
// 		0,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	return client.channel.Consume(
// 		client.queueName,
// 		"",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// }

// // Close will cleanly shut down the channel and connection.
// func (client *Client) Close() error {
// 	if !client.isReady {
// 		return errAlreadyClosed
// 	}
// 	close(client.done)
// 	err := client.channel.Close()
// 	if err != nil {
// 		return error.Error("aaa")
// 	}
// 	err = client.connection.Close()
// 	if err != nil {
// 		return err
// 	}

// 	client.isReady = false
// 	return nil
// }

// ------------------------------ new things ------------------------------
