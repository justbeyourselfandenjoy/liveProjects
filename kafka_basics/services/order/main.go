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
			"bootstrap.servers":         _appCfg.Get("broker_connection", "host") + ":" + _appCfg.Get("broker_connection", "port"),
			"debug":                     _appCfg.Get("broker_connection", "debug"),
			"acks":                      _appCfg.Get("broker_connection", "acks"),
			"socket.timeout.ms":         _appCfg.Get("broker_connection", "socketTimeout"),
			"message.timeout.ms":        _appCfg.Get("broker_connection", "messageTimeout"),
			"go.delivery.report.fields": _appCfg.Get("broker_connection", "goDeliveryReportFields"),
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

	appCfg := newAppConfig()

	//read configuration parameters
	if err := appCfg.initConfigEnv(appCfg.GetValue("service", "name")); err != nil {
		log.Println("Error reading config arguments from env: ", err.Error())
	}

	if err := appCfg.initConfigCL(); err != nil {
		log.Println("Error reading config arguments from command line: ", err.Error())
	}

	var err error
	useZK := appCfg.GetToggle("service", "useZk")
	var zkInstance *zk.Conn

	if useZK {
		if zkInstance, err = appCfg.initConfigZK(); err != nil {
			log.Println("initConfigZK call returned an error: ", err.Error())
			useZK = false
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

	APISchema, err = os.ReadFile(appCfg.Get("api", "APISchemaFile"))
	if err != nil {
		log.Panicln(err)
	}

	//start listening to the configuratiuon chhanges
	if useZK && appCfg.GetToggle("service", "zkHotReload") {
		appCfg.hotReloadConfigZK(zkInstance, []string{"broker_connection", "broker", "api"})
	}

	log.Println("Starting server at " + appCfg.Get("service", "host") + ":" + appCfg.Get("service", "port") + " ... ")
	log.Panicln(http.ListenAndServe(appCfg.Get("service", "host")+":"+appCfg.Get("service", "port"), mux))
}
