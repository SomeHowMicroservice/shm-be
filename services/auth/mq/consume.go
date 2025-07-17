package mq

import (
	"log"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func ConsumeMessage(ch *amqp091.Channel, queueName string, handler func([]byte) error) error {
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.Qos(5, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		go func(workerID int) {
			for msg := range msgs {
				processWithRetry(msg.Body, handler, workerID)
			}
		}(i)
	}

	return nil
}

func processWithRetry(body []byte, handler func([]byte) error, workerID int) {
	maxAttempts := 5
	initialInterval := 1000 * time.Millisecond
	multiplier := 2.0
	maxInterval := 10000 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := handler(body)
		if err == nil {
			return 
		}
		log.Printf("Worker %d: Attempt %d/%d failed: %v", workerID, attempt, maxAttempts, err)

		if attempt < maxAttempts {
			delay := float64(initialInterval) * math.Pow(multiplier, float64(attempt-1))
			if delay > float64(maxInterval) {
				delay = float64(maxInterval)
			}
			time.Sleep(time.Duration(delay))
		}
	}
	log.Printf("Worker %d: Message failed after %d attempts", workerID, maxAttempts)
}
