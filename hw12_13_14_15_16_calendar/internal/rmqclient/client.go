package rmqclient

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewRabbitMQClient(rabbitURL string) (*RabbitMQClient, error) {
	client := &RabbitMQClient{url: rabbitURL}

	var err error
	for i := 0; i < 2; i++ {
		client.conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Error connection to RabbitMq %d: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	client.channel, err = client.conn.Channel()
	if err != nil {
		client.conn.Close()
		return nil, err
	}

	return client, nil
}

func (c *RabbitMQClient) DeclareQueue(queueName string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,
	)
}

func (c *RabbitMQClient) PublishMessage(queueName string, body []byte) error {
	return c.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *RabbitMQClient) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queueName,
		"",
		true,  // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,
	)
}

func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
