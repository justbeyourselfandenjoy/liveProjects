package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-zookeeper/zk"
)

type AppConfig struct {
	groups []ConfigGroup
}

type ConfigGroup struct {
	name     string
	priority uint
	items    []ConfigItem
}

type ConfigItem struct {
	name         string
	value        string
	defaultValue string
	alias        string
}

var appCfg = AppConfig{
	[]ConfigGroup{
		{
			name:     "service",
			priority: 0,
			items: []ConfigItem{
				{name: "name", value: "Order", defaultValue: "Order", alias: "service"},
				{name: "host", value: "localhost", defaultValue: "localhost", alias: "host"},
				{name: "port", value: "8080", defaultValue: "8080", alias: "port"},
			},
		},
		{
			name:     "broker_connection",
			priority: 1,
			items: []ConfigItem{
				{name: "host", value: "localhost", defaultValue: "localhost", alias: "broker_host"},
				{name: "port", value: "9092", defaultValue: "9092", alias: "broker_port"},
				{name: "debug", value: "generic,broker,topic,msg", defaultValue: "generic", alias: "broker_debug"},
				{name: "acks", value: "all", defaultValue: "all", alias: "broker_acks"},
				{name: "produceTimeout", value: "500", defaultValue: "500", alias: "broker_produce_timeout"},
				{name: "socketTimeout", value: "10", defaultValue: "10", alias: "broker_socket_timeout"},
				{name: "messageTimeout", value: "10", defaultValue: "10", alias: "broker_message_timeout"},
				{name: "goDeliveryReportFields", value: "key,value,headers", defaultValue: "all", alias: "broker_go_delivery_report_fields"},
			},
		},
		{
			name:     "broker",
			priority: 2,
			items: []ConfigItem{
				{name: "topicName", value: "OrderReceived", defaultValue: "OrderReceived", alias: "broker_topic_name"},
				{name: "eventName", value: "OrderReceivedEvent", defaultValue: "OrderReceivedEvent", alias: "broker_event_name"},
			},
		},
		{
			name:     "zk",
			priority: 3,
			items: []ConfigItem{
				{name: "host", value: "127.0.0.1", defaultValue: "127.0.0.1", alias: "zk_host"},
				{name: "sessionTimeout", value: "1", defaultValue: "1", alias: "zk_session_timeout"},
				{name: "rootNode", value: "/kafka_basics/Order", defaultValue: "/kafka_basics/Order", alias: "zk_root_node"},
			},
		},
		{
			name:     "api",
			priority: 4,
			items: []ConfigItem{
				{name: "APISchemaFile", value: "api/swagger.yml", defaultValue: "api/swagger.yml", alias: "api_schema_file"},
			},
		},
	},
}

const (
	serviceName       = "Order"
	defaultServerHost = "localhost"
	defaultServerPort = "8080"

	defaultBrokerHost                   = "localhost"
	defaultBrokerPort                   = "9092"
	defaultBrokerDebug                  = "generic" //"generic,broker,topic,msg"
	defaultBrokerAcks                   = "all"
	defaultBrokerProduceTimeout         = 500 //ms
	defaultBrokerSocketTimeout          = 10
	defaultBrokerMessageTimeout         = 10
	defaultBrokerGoDeliveryReportFields = "all" //"key,value,headers"

	defaultOrderReceivedTopicName = "OrderReceived"
	defaultOrderReceivedEventName = "OrderReceivedEvent"

	defaultAPISchemaFile = "api/swagger.yml"

	defaultZKHost           = "127.0.0.1"
	defaultZKSessionTimeout = time.Second
	defaultZKRootNode       = "/kafka_basics/" + serviceName

// defaultZKArguments      = "kafka_topic, kafka_event, api_schema_file"
)

var (
	serverHost = defaultServerHost
	serverPort = defaultServerPort

	brokerHost                   = defaultBrokerHost
	brokerPort                   = defaultBrokerPort
	brokerDebug                  = defaultBrokerDebug
	brokerAcks                   = defaultBrokerAcks
	brokerProduceTimeout         = defaultBrokerProduceTimeout
	brokerSocketTimeout          = defaultBrokerSocketTimeout
	brokerMessageTimeout         = defaultBrokerMessageTimeout
	brokerGoDeliveryReportFields = defaultBrokerGoDeliveryReportFields

	OrderReceivedTopicName = defaultOrderReceivedTopicName
	OrderReceivedEventName = defaultOrderReceivedEventName

	APISchemaFile = defaultAPISchemaFile

	zkHost           = defaultZKHost
	zkSessionTimeout = defaultZKSessionTimeout
	zkRootNode       = defaultZKRootNode

// zkArguments      = defaultZKArguments
)

// read env variables
func initConfigEnv() error {
	if len(os.Getenv("host")) != 0 {
		serverHost = os.Getenv("host")
	}
	if len(os.Getenv("port")) != 0 {
		serverPort = os.Getenv("port")
	}
	return nil
}

func initConfigCL() (bool, bool, error) {
	//override defaults and env variables with application arguments from the CLI
	var serverHostCLArg, serverPortCLArg string
	var useZK, zkHotReload bool
	flag.StringVar(&serverHostCLArg, "host", defaultServerHost, "a server hostname to run the service") //overrides even if it is set by Getenv
	flag.StringVar(&serverPortCLArg, "port", defaultServerPort, "a server port to run the service")
	flag.BoolVar(&useZK, "zk", true, "enable using zookeeper for reading configuration parameters")
	flag.BoolVar(&zkHotReload, "zk_hot_reload", true, "enable re-reading configuration parameters from zk w/o restarting the app")
	flag.Parse()
	// to prevent overriding with the default CLI arg even if it is not set
	flag.Visit(checkCLArgIsSpecified)
	return useZK, zkHotReload, nil
}

func checkCLArgIsSpecified(flag *flag.Flag) {
	if flag.Name == "host" {
		serverHost = flag.Value.String()
	}
	if flag.Name == "port" {
		serverPort = flag.Value.String()
	}
}

func initConfigZK() (*zk.Conn, error) {
	//read from zk
	zkInstance, _, err := zk.Connect([]string{zkHost}, zkSessionTimeout)
	if err != nil {
		log.Panicln(err)
	}
	//checking if the root node exists on zk
	if ok, _, err := zkInstance.Exists(zkRootNode); !ok {
		log.Printf("Can't read the root node [%s], skipping zk\n", zkRootNode)
		return zkInstance, err
	}

	//reading the parameters one-by-one
	zkArg := zkRootNode + "/broker_host"
	if argValue, err := zkGetArg(zkInstance, zkArg); err == nil {
		log.Printf("Read from %s: %s\n", zkArg, string(argValue))
		brokerHost = argValue
	} else {
		log.Printf("Can't read [%v]: %s. Skipping ...\n", zkArg, err.Error())
	}
	//TODO: cycle
	zkArg = zkRootNode + "/broker_port"
	if argValue, err := zkGetArg(zkInstance, zkArg); err == nil {
		log.Printf("Read from %s: %s\n", zkArg, string(argValue))
		brokerPort = argValue
	} else {
		log.Printf("Can't read [%v]: %s. Skipping ...\n", zkArg, err.Error())
	}

	/*
		go func() {
			//		ev := <-zkEventChannel
			//		log.Printf("Get zkEventChannel #3: %+v\n", ev)

			for event := range zkEventChannel {
				log.Printf("£vent range: %v\n", event)
			}
			log.Println("Exiting!!!!!!!")
		}()
	*/
	//	<-zkEventChannel
	//ev := <-zkEventChannel
	//log.Printf("Get zkEventChannel #1: %+v\n", ev)
	/*
			go func() {
		select {
		case ev := <-zkEventChannel:
			log.Printf("Got zkEventChannel event: %+v\n", ev)
			if ev.Err != nil {
				log.Printf("GetW watcher error: %+v\n", ev.Err)
			}
			if ev.State != 100 && ev.State != 101 {
				log.Printf("Watcher is in [%s] state. [StateConnected] or [StateHasSession] is expected\n", ev.State)
				return
			}
		case <-time.After(2 * time.Second):
			log.Println("zx GetW watcher timed out")
		default:
			log.Println("Got default zk event")
		}

			}()
	*/

	return zkInstance, nil
}

func hotReloadConfigZK(zkInstance *zk.Conn) {

	kafkaReloadSig := make(chan bool)

	//TODO: cycle
	go func() {
		zkArg := zkRootNode + "/broker_host"
		for {
			if argValue, err := zkGetArgW(zkInstance, zkArg); err == nil {
				log.Printf("Got zk node change event: %s\n", argValue)
				brokerHost = argValue
				kafkaReloadSig <- true
			} else {
				log.Printf("Can't read [%v]: %s. Skipping...\n", zkArg, err.Error())
			}
		}
	}()

	go func() {
		zkArg := zkRootNode + "/broker_port"
		for {
			if argValue, err := zkGetArgW(zkInstance, zkArg); err == nil {
				log.Printf("Got zk node change event: %s\n", argValue)
				brokerPort = argValue
				kafkaReloadSig <- true
			} else {
				log.Printf("Can't read [%v]: %s. Skipping...\n", zkArg, err.Error())
			}
		}
	}()

	go func() {
		var err error
		for {
			if <-kafkaReloadSig {
				kafkaProducer.Close()
				if kafkaProducer, err = createKafkaProducerInstance(); err != nil {
					log.Panicln(err)
				}
			}
		}
	}()
}

func zkGetArg(zkInstance *zk.Conn, zkArg string) (string, error) {
	argValue, _, err := zkInstance.Get(zkArg)
	if err != nil {
		return "", err
	}
	return string(argValue), nil
}

func zkGetArgW(zkInstance *zk.Conn, zkArg string) (string, error) {

	var argValue []byte

	argValue, _, zkEventChannel, err := zkInstance.GetW(zkArg)
	if err != nil {
		return "", err
	}
	for event := range zkEventChannel {
		log.Printf("£vent range: %v\n", event)
		switch event.Type {
		case 1: //EventNodeCreated
			log.Println("EventNodeCreated event")
		case 2: //EventNodeDeleted
			log.Println("EventNodeDeleted event")
		case 3: //EventNodeDataChanged
			log.Println("EventNodeDataChanged event")
			if argValueTmp, err := zkGetArg(zkInstance, zkArg); err == nil {
				log.Printf("Read from %s: %s\n", zkArg, string(argValueTmp))
				argValue = []byte(argValueTmp)
			} else {
				log.Printf("Can't read [%v]: %s. Skipping ...\n", zkArg, err.Error())
			}
		case 4: //EventNodeChildrenChanged
			log.Println("EventNodeChildrenChanged event")
		case -1: //EventSession
			log.Println("EventSession event")
		case -2: //EventNotWatching
			log.Println("EventNotWatching event")
		}
	}
	return string(argValue), nil
}
