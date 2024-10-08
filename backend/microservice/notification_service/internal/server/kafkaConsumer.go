package server

import (
	"log"
	"notifications/internal/config"
	"notifications/internal/services"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	Consumer *kafka.Consumer
}

// NewConsumer initializes a new Kafka consumer with necessary configurations
func NewConsumer(cfg config.KafkaConfig) (*Consumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers":       cfg.BootstrapServers,
		"group.id":                cfg.GroupID,
		"auto.offset.reset":       cfg.AutoOffsetReset,
		"enable.auto.commit":      cfg.EnableAutoCommit,
		"auto.commit.interval.ms": cfg.AutoCommitIntervalMs,
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer: consumer,
	}, nil
}

// Subscriber subscribes to a list of Kafka topics
func (c *Consumer) Subscriber(topics []string) error {
	if err := c.Consumer.SubscribeTopics(topics, nil); err != nil {
		log.Printf("Err subs: %v", err)
		return err
	}

	return nil
}

// Reader continuously reads messages from the subscribed topics
func (c *Consumer) Reader() error {
	defer c.Consumer.Close()
	for {
		msg, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("read err: %v", err)
			continue
		}

		switch *msg.TopicPartition.Topic {
		case "friends_operation":
			go services.SendNotification(msg)
		default:
			log.Printf("Получено сообщение из неизвестного топика %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
		}
	}
}
