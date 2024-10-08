package server

import (
	"InterestingChats/backend/microservice/db/internal/config"
	"InterestingChats/backend/microservice/db/internal/db"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer to handle Kafka message consumption and interactions with the database
type Consumer struct {
	Consumer  *kafka.Consumer
	DBService db.DBService
}

// NewConsumer to create a new Kafka consumer with the provided configuration
func NewConsumer(dbService db.DBService, kafkaConfig config.KafkaConfig) (*Consumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers":       kafkaConfig.BootstrapServers,
		"group.id":                kafkaConfig.GroupID,
		"auto.offset.reset":       kafkaConfig.AutoOffsetReset,
		"enable.auto.commit":      kafkaConfig.EnableAutoCommit,
		"auto.commit.interval.ms": kafkaConfig.AutoCommitIntervalMs,
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer:  consumer,
		DBService: dbService,
	}, nil
}

// Subscriber to subscribe the Kafka consumer to specified topics
func (c *Consumer) Subscriber(topics []string) error {
	if err := c.Consumer.SubscribeTopics(topics, nil); err != nil {
		log.Printf("Err subs: %v", err)
		return err
	}

	return nil
}

// Reader to continuously read messages from Kafka topics and process them
func (c *Consumer) Reader() error {
	defer c.Consumer.Close()
	for {
		msg, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("read err: %v", err)
			continue
		}

		switch *msg.TopicPartition.Topic {
		case "message_operation":
			go c.DBService.SaveMessage(msg)
		default:
			log.Printf("Recived message from unknow topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
		}
	}
}
