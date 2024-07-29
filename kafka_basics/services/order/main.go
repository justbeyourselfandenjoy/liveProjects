package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-zookeeper/zk"
)

var APISchema []byte
var kafkaProducer *kafka.Producer

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /order/{$}", orderHandler)
}

func createKafkaProducerInstance() (*kafka.Producer, error) {
	return kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers":         brokerHost + ":" + brokerPort,
			"debug":                     brokerDebug,
			"acks":                      brokerAcks,
			"socket.timeout.ms":         brokerSocketTimeout,
			"message.timeout.ms":        brokerMessageTimeout,
			"go.delivery.report.fields": brokerGoDeliveryReportFields,
		},
	)
}

func main() {
	startTime := time.Now()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Got %v signal. Exiting. Uptime %s\n", sig.String(), time.Since(startTime).String())
		os.Exit(0)
	}()

	//read configuration parameters
	if err := initConfigEnv(); err != nil {
		log.Println("Can't read config arguments from env: ", err.Error())
	}
	var err error
	useZK, zkHotReload := false, false
	var zkInstance *zk.Conn
	if useZK, zkHotReload, err = initConfigCL(); err != nil {
		log.Println("Can't read command line arguments: ", err.Error())
	}
	if useZK {
		zkInstance, err = initConfigZK()
		if err != nil {
			log.Println("initConfigZK call returned an error: ", err.Error())
		}
		defer zkInstance.Close()
	}

	mux := http.NewServeMux()
	registerHandlers(mux)

	//create kafka producer instance (a global var), will be used for hot reload with new parameters
	if kafkaProducer, err = createKafkaProducerInstance(); err != nil {
		log.Panicln(err)
	}
	defer kafkaProducer.Close()

	APISchema, err = os.ReadFile(APISchemaFile)
	if err != nil {
		log.Panicln(err)
	}

	//start listening to the configuratiuon chhanges
	if useZK && zkHotReload {
		hotReloadConfigZK(zkInstance)
	}

	log.Println("Starting server at " + serverHost + ":" + serverPort + " ... ")
	log.Panicln(http.ListenAndServe(serverHost+":"+serverPort, mux))
}
