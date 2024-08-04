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

	"justbeyourselfandenjoy/kafka_basics/helpers/config_helpers"
	"justbeyourselfandenjoy/kafka_basics/helpers/kafka_helpers"
)

var APISchema []byte
var kafkaProducer *kafka.Producer

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /order/{$}", orderHandler)
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

	appCfg := config_helpers.NewAppConfig(&_appCfg)

	//read configuration parameters
	if err := appCfg.InitConfigEnv(appCfg.GetValue("service", "name")); err != nil {
		log.Println("Error reading config arguments from env: ", err.Error())
	}

	if err := appCfg.InitConfigCL(); err != nil {
		log.Println("Error reading config arguments from command line: ", err.Error())
	}

	var err error
	useZK := appCfg.GetToggle("service", "useZk")
	var zkInstance *zk.Conn

	if useZK {
		if zkInstance, err = appCfg.InitConfigZK(); err != nil {
			log.Println("initConfigZK call returned an error: ", err.Error())
			useZK = false
		}
		defer zkInstance.Close()
	}

	mux := http.NewServeMux()
	registerHandlers(mux)

	//create kafka producer instance (a global var), will be used for hot reload with new parameters
	var kafkaConfigMap = kafka.ConfigMap{
		"bootstrap.servers":         appCfg.Get("broker_connection", "host") + ":" + appCfg.Get("broker_connection", "port"),
		"debug":                     appCfg.Get("broker_connection", "debug"),
		"acks":                      appCfg.Get("broker_connection", "acks"),
		"socket.timeout.ms":         appCfg.Get("broker_connection", "socketTimeout"),
		"message.timeout.ms":        appCfg.Get("broker_connection", "messageTimeout"),
		"go.delivery.report.fields": appCfg.Get("broker_connection", "goDeliveryReportFields"),
	}
	if kafkaProducer, err = kafka_helpers.CreateKafkaProducerInstance(&kafkaConfigMap); err != nil {
		log.Panicln(err)
	}
	defer kafkaProducer.Close()

	APISchema, err = os.ReadFile(appCfg.Get("api", "APISchemaFile"))
	if err != nil {
		log.Panicln(err)
	}

	//start listening to the configuratiuon chhanges
	if useZK && appCfg.GetToggle("service", "zkHotReload") {
		appCfg.HotReloadConfigZK(zkInstance, []string{"broker_connection", "broker", "api"}, kafkaProducer)
	}

	log.Println("Starting server at " + appCfg.Get("service", "host") + ":" + appCfg.Get("service", "port") + " ... ")
	log.Panicln(http.ListenAndServe(appCfg.Get("service", "host")+":"+appCfg.Get("service", "port"), mux))
}
