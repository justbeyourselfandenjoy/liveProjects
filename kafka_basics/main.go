package main

import (
	"justbeyourselfandenjoy/kafka_basics/helpers"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

func main() {
	err := helpers.PublishEvent(
		&kafka.ConfigMap{
			"bootstrap.servers": brokerIP + ":" + brokerPort,
			"debug":             brokerDebug,
			"client.id":         clientID,
			"acks":              brokerAcks},
		"OrderReceived",
		&helpers.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
			EventName:      "OrderReceived",
			EventBody:      "Test message from the app #2",
		},
	)

	if err != nil {
		log.Panicln(err)
	}

	log.Println("Event published successfully")

}
