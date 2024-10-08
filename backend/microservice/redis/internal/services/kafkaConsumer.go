package services

import (
	"InterestingChats/backend/microservice/redis/internal/config"
	"InterestingChats/backend/microservice/redis/internal/consts"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer is a struct that represents a Kafka consumer with an associated Redis client.
type Consumer struct {
	Consumer *kafka.Consumer
	RC       rdb.RedisClient
}

// NewConsumer creates a new Kafka consumer with the provided Redis client.
func NewConsumer(rdbClient *rdb.RedisClient, cfg *config.Config) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       cfg.Kafka.BootstrapServers,
		"group.id":                cfg.Kafka.GroupID,
		"auto.offset.reset":       cfg.Kafka.AutoOffsetReset,
		"enable.auto.commit":      cfg.Kafka.Consumer.EnableAutoCommit,
		"auto.commit.interval.ms": cfg.Kafka.Consumer.AutoCommitIntervalMs,
		"session.timeout.ms":      cfg.Kafka.Consumer.SessionTimeoutMs,
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer: consumer,
		RC:       *rdbClient,
	}, nil
}

// Subscriber subscribes the consumer to the specified topics.
func (c *Consumer) Subscriber(topics []string) error {
	if err := c.Consumer.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	return nil
}

// Reader continuously reads messages from the subscribed topics and processes them.
func (c *Consumer) Reader(cfg *config.Config) error {
	log.Printf("start read topics")
	defer c.Consumer.Close()
	for {
		data, err := c.Consumer.ReadMessage(-1)
		if err != nil {
			continue
		}

		switch *data.TopicPartition.Topic {
		case "session_SET":
			go SetSession(data, &c.RC, cfg.Redis.SessionTTL)
		case "push_friends":
			go AddFriend(data, &c.RC)
		case "push_subscribers":
			go AddSubscribers(data, &c.RC)
		case "update_subscriber":
			go RemoveFriendAndAddSubs(data, &c.RC)
		case "update_friend":
			go RemoveSubsAndAddFriend(data, &c.RC)
		case "session_UPDATE":
			log.Printf("recived")
			go UpdateSessionData(data, &c.RC, cfg.Redis.SessionTTL)
		default:
			log.Printf(consts.InternalUnknowTopic, *data.TopicPartition.Topic, string(data.Value))
		}
	}
}
