package queue

import (
	"time"
	"wall/utils/logger"

	"github.com/streadway/amqp"
)

type queue struct {
	logger logger.Logger

	url  string
	name string

	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	closed       bool
}

func NewQueue(logger logger.Logger, url string, qName string) *queue {
	q := new(queue)
	q.logger = logger
	q.url = url
	q.name = qName

	q.connect()
	go q.reconnector()

	return q
}

func (q *queue) Send(message string) {
	err := q.channel.Publish(
		q.name, // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		logger.Log(q.logger, "Error", "NA", logger.Trace(), err.Error())
	}
}

func (q *queue) Close() {
	logger.Log(q.logger, "Info", "NA", logger.Trace(), "MQ close")
	q.closed = true
	q.channel.Close()
	q.connection.Close()
}

func (q *queue) connect() {
	for {
		logger.Log(q.logger, "Info", "NA", logger.Trace(), "Connecting to rabbitmq on: "+q.url)
		conn, err := amqp.Dial(q.url)
		if err == nil {
			q.connection = conn
			q.errorChannel = make(chan *amqp.Error)
			q.connection.NotifyClose(q.errorChannel)

			logger.Log(q.logger, "Info", "NA", logger.Trace(), "Connection established!")

			q.openChannel()

			return
		}

		logger.Log(q.logger, "Error", "NA", logger.Trace(), "Connection to rabbitmq failed. Retrying in 1 sec... "+err.Error())
		time.Sleep(1000 * time.Millisecond)
	}
}

func (q *queue) openChannel() {
	channel, err := q.connection.Channel()
	if err != nil {
		logger.Log(q.logger, "Error", "NA", logger.Trace(), "Opening channel failed"+err.Error())
	}

	q.channel = channel
}

func (q *queue) reconnector() {
	for {
		err := <-q.errorChannel
		if !q.closed {
			logger.Log(q.logger, "Error", "NA", logger.Trace(), "Reconnecting after connection closed"+err.Error())
			q.connect()
		}
	}
}
