package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"justbeyourselfandenjoy/kafka_basics/helpers/kafka_helpers"
	"justbeyourselfandenjoy/service_order/swagger"
	"log"
	"net/http"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("versionHandler has been called")
	w.Write([]byte("1.0.0"))
}

func upTimeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("upTimeHandler has been called")
	w.Write([]byte(time.Since(appStartTime).String()))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("healthHandler has been called")
	w.WriteHeader(http.StatusOK)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("orderHandler has been called")

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Println("Error reading request body: ", err.Error())
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content-type 'application/json' is allowed, got '%s'", contentType)
		log.Println("Disallowed content type received: ", contentType)
		return
	}

	//TODO to implement the real validation
	if err = validateJsonAgainstSchema(APISchema, bodyBytes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("order JSON against schema validation failed: " + err.Error()))
		log.Println("Error validating the request against schema: ", err.Error())
		return
	}

	var orderReceived swagger.Order

	if err = json.Unmarshal(bodyBytes, &orderReceived); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("Error converting body to Order object: ", err.Error())
		return
	}

	if err = validateOrderPayload(orderReceived); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("order payload validation failed: " + err.Error()))
		log.Println("Invalid payload received: ", err.Error())
		return
	}

	if err := kafkaProducer.PublishEvent(
		kafka_helpers.BuildBaseEvent(_appCfg.Get("broker", "eventNameProduce"), string(bodyBytes)),
		int(_appCfg.GetInt("broker_connection", "produceTimeout")),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Println("Error publishing event to Kafka: ", err.Error())
		return
	}
}
