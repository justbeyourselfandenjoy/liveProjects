package helpers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

// BaseEvent represents common properties of an event
type BaseEvent struct {
	EventID        uuid.UUID
	EventTimestamp time.Time
	EventName      string
	EventBody      string
}

func PublishEvent(kafkaConfig *kafka.ConfigMap, topic string, event *BaseEvent) error {
	log.Println("PublishEvent is called")

	kafkaProducer, err := kafka.NewProducer(kafkaConfig)

	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	var value []byte
	if value, err = json.Marshal(event); err != nil {
		return err
	}

	err = kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value)},
		nil, // delivery channel
	)

	if err != nil {
		return err
	}
	return nil
}
