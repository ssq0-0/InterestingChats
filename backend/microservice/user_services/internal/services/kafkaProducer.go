package services

import (
	"InterestingChats/backend/user_services/internal/config"
	"InterestingChats/backend/user_services/internal/consts"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Producer *kafka.Producer
	Topics   config.ConfigKafkaTopics
}

func NewProducer(bootstrapServers string, topics config.ConfigKafkaTopics) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: producer,
		Topics:   topics,
	}, nil
}

func (p *Producer) Writer(message interface{}, typeMessage int) error {
	log.Printf("записываю в кафку, topic: %v", typeMessage)
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	var topic string
	switch typeMessage {
	case consts.KAFKA_Friendship:
		topic = p.Topics.FriendsOperation
	case consts.KAFKA_Session:
		topic = p.Topics.SessionSet
	case consts.KAFKA_PushFriends:
		topic = p.Topics.PushFriends
	case consts.KAFKA_PushSubscribers:
		topic = p.Topics.PushSubscribers
	case consts.KAFKA_UpdateSubscribers:
		topic = p.Topics.UpdateSubscriber
	case consts.KAFKA_RemoveSubscriberAndAddFriend:
		topic = p.Topics.UpdateFriend
	case consts.KAFKA_RemoveFriendAndAddSubscriber:
		topic = p.Topics.UpdateSubscriber
	case consts.KAFKA_session_UPDATE:
		topic = p.Topics.SessionUpdate
	}

	if err := p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil); err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() {
	p.Producer.Close()
}
