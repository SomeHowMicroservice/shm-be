package mq

import "github.com/rabbitmq/amqp091-go"

func ConsumeMessage(ch *amqp091.Channel, queueName string, handler func([]byte)) error {
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()
	return nil
}