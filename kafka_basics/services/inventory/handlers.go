package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"justbeyourselfandenjoy/kafka_basics/helpers/kafka_helpers"

	"justbeyourselfandenjoy/service_order/swagger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthHandler has been called")
	w.WriteHeader(http.StatusOK)
}

func kafkaMessageHandler(msg *kafka.Message) {

	var baseEvent kafka_helpers.BaseEvent
	useDLQ := _appCfg.GetToggle("service", "useDLQ")

	if msg == nil {
		return
	}

	if err := json.Unmarshal(msg.Value, &baseEvent); err != nil {
		log.Printf("Error converting message to Order object: %s. Skipping the message ...\n", err.Error())
		if useDLQ {
			if err := publishDLQ(string(msg.Value)); err != nil {
				log.Printf("Error creating DLQ: %s\n", err.Error())
			}
		}
		return
	}

	//Checking for the duplicated events
	log.Printf("Checking if the event %s has been processed already ... ", baseEvent.EventID)
	if !eventsRegistry.Exists(baseEvent.EventID) {
		log.Println("New event. Adding to the registry ...")
		eventsRegistry.Add(baseEvent.EventID)
	} else {
		log.Println("Duplicate found. Skipping ...")
		return
	}

	var orderReceived swagger.Order
	if err := json.Unmarshal([]byte(baseEvent.EventBody), &orderReceived); err != nil {
		log.Printf("Error fetching Order object from Kafka event: %s. Skipping the message ...\n", err.Error())
		if useDLQ {
			if err := publishDLQ(string(msg.Value)); err != nil {
				log.Printf("Error creating DLQ: %s\n", err.Error())
			}
		}
		return
	}
	eventsRegistry.Set(baseEvent.EventID, kafka_helpers.EVENT_STATUS_PROCESSING)

	//TODO check for the repeated orders
	log.Printf("Checking if the order %s has been processed already ... done\n", orderReceived.ID)

	//TODO process order
	log.Printf("Decreasing inventory for the items in order %v\n", orderReceived.ID)

	if len(orderReceived.Products) > 0 {
		for _, product := range orderReceived.Products {
			log.Printf("Reducing number of [%s, %s] by %v\n", product.ID, product.ProductCode, product.Quantity)
		}
	} else {
		log.Printf("[WARN] Got an order with empty products. Skipping it...")
		if useDLQ {
			if err := publishDLQ(string(msg.Value)); err != nil {
				log.Printf("Error creating DLQ: %s\n", err.Error())
			}
		}
		eventsRegistry.Set(baseEvent.EventID, kafka_helpers.EVENT_STATUS_PROCESSING_FAILED)
		return
	}

	//sending the event about the processing completion
	if err := kafkaProducer.PublishEvent(
		kafka_helpers.BuildBaseEvent(_appCfg.Get("broker", "eventNameProduce"), baseEvent.EventBody),
		int(_appCfg.GetInt("broker_connection", "produceTimeout")),
	); err != nil {
		if useDLQ {
			if err := publishDLQ(string(msg.Value)); err != nil {
				log.Printf("Error creating DLQ: %s\n", err.Error())
			}
		}
		log.Println("Error publishing event to Kafka: ", err.Error())
		eventsRegistry.Set(baseEvent.EventID, kafka_helpers.EVENT_STATUS_PROCESSING_FAILED)
		return
	}
	eventsRegistry.Set(baseEvent.EventID, kafka_helpers.EVENT_STATUS_PROCESSED)
	log.Printf("**************** Events registry state: %v\n", eventsRegistry.String())
}

func publishDLQ(msg string) error {
	if err := kafkaProducerDLQ.PublishEvent(
		&kafka_helpers.BaseEvent{
			EventID:        uuid.New(),
			EventTimestamp: time.Now(),
			EventName:      _appCfg.Get("broker", "eventNameDLQ"),
			EventBody:      msg,
		},
		int(_appCfg.GetInt("broker_connection", "produceTimeout")),
	); err != nil {
		log.Println("Error publishing event to Kafka: ", err.Error())
		return err
	}
	return nil
}
