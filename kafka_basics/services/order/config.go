package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

var _appCfg = AppConfig{
	groups: map[string]ConfigGroup{
		"service": {
			priority: 1,
			items: map[string]ConfigItem{
				"name":        {value: "Order", defaultValue: "Order", alias: "service_name", description: "service name, internal, used for env variables prefix"},
				"host":        {value: "localhost", defaultValue: "localhost", alias: "service_host", description: "host to bind the server"},
				"port":        {value: "8080", defaultValue: "8080", alias: "service_port", description: "port to bind the server"},
				"useZk":       {value: "true", defaultValue: "true", alias: "use_zk", description: "toggle for using zookeeper"},
				"zkHotReload": {value: "true", defaultValue: "true", alias: "zk_hot_reload", description: "toggle to enable hot reload for parameters from zookeeper"},
			},
		},
		"broker_connection": {
			priority: 2,
			items: map[string]ConfigItem{
				"host":                   {value: "localhost", defaultValue: "localhost", alias: "broker_host", description: "broker host to connect"},
				"port":                   {value: "9092", defaultValue: "9092", alias: "broker_port", description: "broker port to connect"},
				"debug":                  {value: "generic,broker,topic,msg", defaultValue: "generic", alias: "broker_debug", description: "broker client debug level"},
				"acks":                   {value: "all", defaultValue: "all", alias: "broker_acks", description: "broker client acks parameter"},
				"produceTimeout":         {value: "500", defaultValue: "500", alias: "broker_produce_timeout", description: "timeout for producing messages to the broker client"},
				"socketTimeout":          {value: "10", defaultValue: "10", alias: "broker_socket_timeout", description: "socket.timeout.ms for broker producer"},
				"messageTimeout":         {value: "10", defaultValue: "10", alias: "broker_message_timeout", description: "message.timeout.ms for broker producer"},
				"goDeliveryReportFields": {value: "key,value,headers", defaultValue: "all", alias: "broker_go_delivery_report_fields", description: "go.delivery.report.fields for broker producer"},
			},
		},
		"broker": {
			priority: 3,
			items: map[string]ConfigItem{
				"topicName": {value: "OrderReceived", defaultValue: "OrderReceived", alias: "broker_topic_name", description: "broker's topic name used by the service"},
				"eventName": {value: "OrderReceivedEvent", defaultValue: "OrderReceivedEvent", alias: "broker_event_name", description: "broker's event name used by the service"},
			},
		},
		"zk": {
			priority: 4,
			items: map[string]ConfigItem{
				"host":           {value: "127.0.0.1", defaultValue: "127.0.0.1", alias: "zk_host", description: "zookeepeer host to connect"},
				"sessionTimeout": {value: "1", defaultValue: "1", alias: "zk_session_timeout", description: "zookeepeer session timeout"},
				"rootNode":       {value: "/kafka_basics/Order", defaultValue: "/kafka_basics/Order", alias: "zk_root_node", description: "zookeepeer root node for the service to read"},
			},
		},
		"api": {
			priority: 4,
			items: map[string]ConfigItem{
				"APISchemaFile": {value: "api/swagger.yml", defaultValue: "api/swagger.yml", alias: "api_schema_file", description: "schema file name to validate requests against"},
			},
		},
	},
}

var _CLFlags map[string]struct{}

type AppConfig struct {
	sync.Mutex
	groups map[string]ConfigGroup
}

type ConfigGroup struct {
	priority uint
	items    map[string]ConfigItem
}

type ConfigItem struct {
	value        string
	defaultValue string
	alias        string
	description  string
}

func newAppConfig() *AppConfig {
	return &_appCfg
}

func (c *AppConfig) Get(group string, name string) string {
	return c.GetValue(group, name)
}

func (c *AppConfig) GetValue(group string, name string) string {
	return c.groups[group].items[name].value
}
func (c *AppConfig) GetAlias(group string, name string) string {
	return c.groups[group].items[name].alias
}

func (c *AppConfig) GetDefault(group string, name string) string {
	return c.groups[group].items[name].defaultValue
}

func (c *AppConfig) GetToggle(group string, name string) bool {
	return c.groups[group].items[name].value == "true"
}

func (c *AppConfig) GetInt(group string, name string) uint64 {
	if s, err := strconv.ParseUint(c.groups[group].items[name].value, 10, 32); err == nil {
		return s
	}
	return 0
}

func (c *AppConfig) Set(group string, name string, value string) error {
	c.Lock()
	item := c.groups[group].items[name]
	item.value = value
	c.groups[group].items[name] = item
	defer c.Unlock()
	return nil
}

// read env variables
func (c *AppConfig) initConfigEnv(prefix string) error {
	for groupName, configGroup := range c.groups {
		for configName, configItem := range configGroup.items {
			envVar := os.Getenv(prefix + "_" + configItem.alias)
			if len(envVar) != 0 {
				log.Printf("setting [%v.%s] to value [%s] from env valiable [%s]\n", groupName, configName, envVar, prefix+"_"+configItem.alias)
				c.Set(groupName, configName, envVar)
			}

		}
	}
	return nil
}

func (c *AppConfig) initConfigCL() error {
	for _, configGroup := range c.groups {
		for _, configItem := range configGroup.items {
			var CLValue string
			flag.StringVar(&CLValue, configItem.alias, configItem.defaultValue, configItem.description)
		}
	}
	flag.Parse()
	_CLFlags = make(map[string]struct{})
	flag.Visit(setCLArgs)

	//checking what flags are set from the CL and set only those
	for groupName, configGroup := range c.groups {
		for configName, configItem := range configGroup.items {
			flagValue := flag.Lookup(configItem.alias).Value.String()
			//override only actually set CL parameters
			if _, ok := _CLFlags[configItem.alias]; ok {
				log.Printf("setting [%v.%s] to new value [%s] from CL argument\n", groupName, configName, flagValue)
				c.Set(groupName, configName, flagValue)
			}
		}
	}
	return nil
}

func setCLArgs(flag *flag.Flag) {
	_CLFlags[flag.Name] = struct{}{}
}

func (c *AppConfig) initConfigZK() (*zk.Conn, error) {
	//read from zk
	zkInstance, _, err := zk.Connect([]string{c.Get("zk", "host")}, time.Second)
	if err != nil {
		return nil, err
	}

	//checking if the root node exists on zk
	zkRootNode := c.Get("zk", "rootNode")
	if ok, _, err := zkInstance.Exists(zkRootNode); !ok {
		log.Printf("initConfigZK: can't read the root node [%s]: %s. Skipping zk...\n", zkRootNode, err.Error())
		return zkInstance, err
	}

	//reading the parameters
	for groupName, configGroup := range c.groups {
		for configName, configItem := range configGroup.items {
			zkArg := zkRootNode + "/" + configItem.alias
			if argValue, err := c.zkGetArg(zkInstance, zkArg); err == nil {
				log.Printf("initConfigZK: read from %s: %s\n", zkArg, string(argValue))
				c.Set(groupName, configName, argValue)
			} else {
				log.Printf("initConfigZK: can't read [%v]: %s. Skipping ...\n", zkArg, err.Error())
			}
		}
	}

	return zkInstance, nil
}

func (c *AppConfig) hotReloadConfigZK(zkInstance *zk.Conn, configGroups []string) {

	kafkaReloadSig := make(chan bool)

	//checking if the root node exists on zk
	zkRootNode := c.Get("zk", "rootNode")
	if ok, _, _ := zkInstance.Exists(zkRootNode); !ok {
		log.Printf("hotReloadConfigZK: can't read the root node [%s], skipping hotReloadConfigZK...\n", zkRootNode)
		return
	}

	for _, configGroup := range configGroups {
		for configName, configItem := range c.groups[configGroup].items {
			go func() {
				zkArg := zkRootNode + "/" + configItem.alias
				for {
					if argValue, err := c.zkGetArgW(zkInstance, zkArg); err == nil {
						log.Printf("hotReloadConfigZK: got zk node change event: %s\n", argValue)
						c.Set(configGroup, configName, argValue)
						kafkaReloadSig <- true
					} else {
						log.Printf("hotReloadConfigZK: can't read [%v]: %s. Skipping...\n", zkArg, err.Error())
						return
					}
				}
			}()
		}
	}

	go func() {
		var err error
		for {
			if <-kafkaReloadSig {
				//TODO wait untile the previous relod is completed
				kafkaProducer.Close()
				if kafkaProducer, err = createKafkaProducerInstance(); err != nil {
					log.Panicln(err)
				}
			}
		}
	}()
}

func (c *AppConfig) zkGetArg(zkInstance *zk.Conn, zkArg string) (string, error) {
	argValue, _, err := zkInstance.Get(zkArg)
	if err != nil {
		return "", err
	}
	return string(argValue), nil
}

func (c *AppConfig) zkGetArgW(zkInstance *zk.Conn, zkArg string) (string, error) {

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
			if argValueTmp, err := c.zkGetArg(zkInstance, zkArg); err == nil {
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
