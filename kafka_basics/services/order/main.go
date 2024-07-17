package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var APISchema []byte
var kafkaProducer *kafka.Producer

func main() {
	startTime := time.Now()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Got %v signal. Exiting. Uptime %v\n", sig.String(), time.Since(startTime).String())
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	var err error
	APISchema, err = os.ReadFile(APISchemaFile)
	if err != nil {
		log.Panicln(err)
	}
	mux.HandleFunc("POST /order/{$}", orderHandler)

	kafkaProducer, err = kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": brokerIP + ":" + brokerPort,
			"debug":             brokerDebug,
			"acks":              brokerAcks},
	)

	if err != nil {
		log.Panicln(err)
	}
	defer kafkaProducer.Close()

	log.Println("Starting server at " + serverIP + ":" + serverPort + " ... ")
	log.Panicln(http.ListenAndServe(serverIP+":"+serverPort, mux))
}
