package kafka_helpers

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

func BuildBaseEvent(name string, body string) *BaseEvent {
	return &BaseEvent{
		EventID:        uuid.New(),
		EventTimestamp: time.Now(),
		EventName:      name,
		EventBody:      body,
	}
}

func PublishEvent(kafkaProducer *kafka.Producer, topic string, event *BaseEvent) error {
	log.Println("PublishEvent is called")

	var value []byte
	var err error

	// TODO implement protobuf
	if value, err = json.Marshal(event); err != nil {
		return err
	}

	//TODO make async call
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
