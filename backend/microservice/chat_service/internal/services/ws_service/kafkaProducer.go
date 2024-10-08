package ws_service

import (
	"chat_service/internal/config"
	"chat_service/internal/models"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer represents a Kafka producer that can send messages to a specified topic.
type Producer struct {
	Producer *kafka.Producer
	Topic    []string
}

// NewProducer initializes a new Kafka producer with the given configuration.
// It returns a pointer to the Producer and any error encountered during initialization.
func NewProducer(kafkaConfig config.KafkaConfig) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.BootstrapServers,
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: producer,
		Topic:    kafkaConfig.Topics,
	}, nil
}

// Writer sends a message to the Kafka topic.
// It marshals the message to JSON and produces it to Kafka.
func (p *Producer) Writer(message *models.Message) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	topic := "message_operation"
	if err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg)},
		nil,
	); err != nil {
		return err
	}

	return nil
}

// Close closes the Kafka producer, flushing any remaining messages.
func (p *Producer) Close() {
	p.Producer.Close()
}
