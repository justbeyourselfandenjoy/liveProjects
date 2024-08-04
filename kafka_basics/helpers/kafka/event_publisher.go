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

func PublishEvent(kafkaProducer *kafka.Producer, topic string, event *BaseEvent, produceTimeoutMs int) error {
	log.Println("PublishEvent is called")

	var value []byte
	var err error

	// TODO implement protobuf
	if value, err = json.Marshal(event); err != nil {
		return err
	}

	// Produce the message asynchronously
	go func() {
		kafkaProducer.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value)}
	}()

	// Delivery report handler for produced messages
	go func() {
		for e := range kafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed producing message to Kafka: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			case kafka.Error:
				if ev.IsFatal() && ev.Code() == kafka.ErrLeaderNotAvailable {
					log.Printf("Kafka error: connection refused: %s\n", ev.String())
				} else {
					log.Printf("Kafka error: error producing message: %s\n", ev.String())
				}
				/*
					case *kafka.Stats:
						var stats map[string]interface{}
						json.Unmarshal([]byte(ev.String()), &stats)
						log.Printf("Got kafka.Stats event. Stats: %v messages (%v bytes) messages consumed %v\n", e, stats["rxmsgs"], stats["rxmsg_bytes"])
				*/
			default:
				log.Printf("Ignored Kafka event: %v\n", e)
			}
		}
	}()

	// Wait for message deliveries before shutting down
	if outstandingNum := kafkaProducer.Flush(produceTimeoutMs); outstandingNum > 0 {
		log.Printf("Expected empty queue after Flush(), still has %d\n", outstandingNum)
	}

	/*
		if err != nil {
			return err
		}
	*/
	return nil
}
