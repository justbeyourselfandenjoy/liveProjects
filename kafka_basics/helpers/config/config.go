package config_helpers

import (
	"flag"
	"fmt"
	"justbeyourselfandenjoy/kafka_basics/helpers/kafka_helpers"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-zookeeper/zk"
)

var _CLFlags map[string]struct{}

type AppConfig struct {
	sync.Mutex
	Groups map[string]ConfigGroup
}

type ConfigGroup struct {
	Priority uint
	Items    map[string]ConfigItem
}

type ConfigItem struct {
	Value        string
	DefaultValue string
	Alias        string
	Description  string
}

func NewAppConfig(config *AppConfig) *AppConfig {
	return config
}

func (c *AppConfig) Get(group string, name string) string {
	return c.GetValue(group, name)
}

func (c *AppConfig) GetValue(group string, name string) string {
	return c.Groups[group].Items[name].Value
}
func (c *AppConfig) GetAlias(group string, name string) string {
	return c.Groups[group].Items[name].Alias
}

func (c *AppConfig) GetDefault(group string, name string) string {
	return c.Groups[group].Items[name].DefaultValue
}

func (c *AppConfig) GetToggle(group string, name string) bool {
	return c.Groups[group].Items[name].Value == "true"
}

func (c *AppConfig) GetInt(group string, name string) uint64 {
	if s, err := strconv.ParseUint(c.Groups[group].Items[name].Value, 10, 32); err == nil {
		return s
	}
	return 0
}

func (c *AppConfig) Set(group string, name string, value string) error {
	c.Lock()
	item := c.Groups[group].Items[name]
	item.Value = value
	c.Groups[group].Items[name] = item
	defer c.Unlock()
	return nil
}

func (c *AppConfig) GetConfigValues() string {
	var _cfg string
	for groupName, configGroup := range c.Groups {
		for configName, configItem := range configGroup.Items {
			_cfg += fmt.Sprintf("%v.%s=%v\n", groupName, configName, configItem.Value)
		}
	}
	return _cfg
}

// read env variables
func (c *AppConfig) InitConfigEnv(prefix string) error {
	for groupName, configGroup := range c.Groups {
		for configName, configItem := range configGroup.Items {
			envVar := os.Getenv(prefix + "_" + configItem.Alias)
			if len(envVar) != 0 {
				log.Printf("setting [%v.%s] to value [%s] from env valiable [%s]\n", groupName, configName, envVar, prefix+"_"+configItem.Alias)
				c.Set(groupName, configName, envVar)
			}

		}
	}
	return nil
}

func (c *AppConfig) InitConfigCL() error {
	for _, configGroup := range c.Groups {
		for _, configItem := range configGroup.Items {
			var CLValue string
			flag.StringVar(&CLValue, configItem.Alias, configItem.DefaultValue, configItem.Description)
		}
	}
	flag.Parse()
	_CLFlags = make(map[string]struct{})
	flag.Visit(setCLArgs)

	//checking what flags are set from the CL and set only those
	for groupName, configGroup := range c.Groups {
		for configName, configItem := range configGroup.Items {
			//override only actually set CL parameters
			if _, ok := _CLFlags[configItem.Alias]; ok {
				flagValue := flag.Lookup(configItem.Alias).Value.String()
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

func (c *AppConfig) InitConfigZK() (*zk.Conn, error) {
	//read from zk
	zkInstance, _, err := zk.Connect([]string{c.Get("zk", "host")}, time.Second)
	if err != nil {
		return nil, err
	}

	//checking if the root node exists on zk
	zkRootNode := c.Get("zk", "rootNode")

	if ok, _, err := zkInstance.Exists(zkRootNode); !ok {
		log.Printf("initConfigZK: can't read the root node [%s]. Skipping zk...\n", zkRootNode)
		return zkInstance, err
	}

	//reading the parameters
	for groupName, configGroup := range c.Groups {
		for configName, configItem := range configGroup.Items {
			zkArg := zkRootNode + "/" + configItem.Alias
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

func (c *AppConfig) HotReloadConfigZK(zkInstance *zk.Conn, configGroups []string, kafkaProducer *kafka.Producer) {

	kafkaReloadSig := make(chan bool)

	//checking if the root node exists on zk
	zkRootNode := c.Get("zk", "rootNode")
	if ok, _, _ := zkInstance.Exists(zkRootNode); !ok {
		log.Printf("hotReloadConfigZK: can't read the root node [%s], skipping hotReloadConfigZK...\n", zkRootNode)
		return
	}

	for _, configGroup := range configGroups {
		for configName, configItem := range c.Groups[configGroup].Items {
			go func() {
				zkArg := zkRootNode + "/" + configItem.Alias
				for {
					if argValue, err := c.zkGetArgW(zkInstance, zkArg); err == nil {
						log.Printf("hotReloadConfigZK: got zk node change event: %s\n", argValue)
						c.Set(configGroup, configName, argValue)
						if configGroup == "broker_connection" {
							kafkaReloadSig <- true
						}
					} else {
						log.Printf("hotReloadConfigZK: can't read [%v]: %s. Skipping...\n", zkArg, err.Error())
						return
					}
				}
			}()
		}
	}

	go func() {
		for {
			if <-kafkaReloadSig {
				//TODO wait until the previous relod is completed
				kafkaProducer.Close()
				var err error
				//TODO rework, too dirty
				kafkaConfigMap := kafka.ConfigMap{
					"bootstrap.servers":         c.Get("broker_connection", "host") + ":" + c.Get("broker_connection", "port"),
					"debug":                     c.Get("broker_connection", "debug"),
					"acks":                      c.Get("broker_connection", "acks"),
					"socket.timeout.ms":         c.Get("broker_connection", "socketTimeout"),
					"message.timeout.ms":        c.Get("broker_connection", "messageTimeout"),
					"go.delivery.report.fields": c.Get("broker_connection", "goDeliveryReportFields"),
				}
				if kafkaProducer, err = kafka_helpers.CreateKafkaProducerInstance(&kafkaConfigMap); err != nil {
					log.Panicln(err)
				}
				log.Print("Updated configuration snapshot \n", c.GetConfigValues())
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
