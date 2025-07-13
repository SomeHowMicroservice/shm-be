package initialization

import (
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/services/auth/config"
	"github.com/rabbitmq/amqp091-go"
)

func InitMessageQueue(cfg *config.Config) (*amqp091.Connection, *amqp091.Channel, error) {
	dsn := fmt.Sprintf("amqps://%s:%s@%s/%s",
		cfg.MessageQueue.RUser,
		cfg.MessageQueue.RPassword,
		cfg.MessageQueue.RHost,
		cfg.MessageQueue.RUser,
	)
	conn, err := amqp091.Dial(dsn)
	if err != nil {
		log.Fatalf("Kết nối RabbitMQ thất bại: %v", err)
	}
	chann, err := conn.Channel()
	if err != nil {
		log.Fatalf("Mở 1 kênh MQ thất bại: %v", err)
	}
	return conn, chann, nil
}
