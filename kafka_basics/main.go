package main

import (
	kafka_helpers "justbeyourselfandenjoy/kafka_basics/helpers"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	kafkaProducer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": brokerIP + ":" + brokerPort,
			"debug":             brokerDebug,
			"acks":              brokerAcks},
	)

	if err != nil {
		log.Panicln(err)
	}
	defer kafkaProducer.Close()

	err = kafka_helpers.PublishEvent(
		kafkaProducer,
		"OrderReceived",
		kafka_helpers.BuildBaseEvent("OrderReceived", "Test message from the app #2"),
		500,
	)

	if err != nil {
		log.Panicln(err)
	}

	log.Println("Event published successfully")

}
