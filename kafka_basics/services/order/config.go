package main

import "justbeyourselfandenjoy/kafka_basics/helpers/config_helpers"

var _appCfg = config_helpers.AppConfig{
	Groups: map[string]config_helpers.ConfigGroup{
		"service": {
			Priority: 1,
			Items: map[string]config_helpers.ConfigItem{
				"name":        {Value: "Order", DefaultValue: "Order", Alias: "service_name", Description: "service name, internal, used for env variables prefix"},
				"host":        {Value: "localhost", DefaultValue: "localhost", Alias: "service_host", Description: "host to bind the server"},
				"port":        {Value: "8080", DefaultValue: "8080", Alias: "service_port", Description: "port to bind the server"},
				"useZk":       {Value: "true", DefaultValue: "true", Alias: "use_zk", Description: "toggle for using zookeeper"},
				"zkHotReload": {Value: "true", DefaultValue: "true", Alias: "zk_hot_reload", Description: "toggle to enable hot reload for parameters from zookeeper"},
			},
		},
		"broker_connection": {
			Priority: 2,
			Items: map[string]config_helpers.ConfigItem{
				"host":  {Value: "localhost", DefaultValue: "localhost", Alias: "broker_host", Description: "broker host to connect"},
				"port":  {Value: "9092", DefaultValue: "9092", Alias: "broker_port", Description: "broker port to connect"},
				"debug": {Value: "generic,broker,topic,msg", DefaultValue: "generic", Alias: "broker_debug", Description: "broker client debug level"},
				"acks":  {Value: "all", DefaultValue: "all", Alias: "broker_acks", Description: "broker client acks parameter"},
				//TODO Configuration property `fetch.wait.max.ms` (500) should be set lower than `socket.timeout.ms` (10) by at least 1000ms to avoid blocking and timing out sub-sequent requests
				"produceTimeout":         {Value: "500", DefaultValue: "500", Alias: "broker_produce_timeout", Description: "timeout for producing messages to the broker client"},
				"socketTimeout":          {Value: "10", DefaultValue: "10", Alias: "broker_socket_timeout", Description: "socket.timeout.ms for broker producer"},
				"messageTimeout":         {Value: "10", DefaultValue: "10", Alias: "broker_message_timeout", Description: "message.timeout.ms for broker producer"},
				"goDeliveryReportFields": {Value: "key,value,headers", DefaultValue: "all", Alias: "broker_go_delivery_report_fields", Description: "go.delivery.report.fields for broker producer"},
			},
		},
		"broker": {
			Priority: 3,
			Items: map[string]config_helpers.ConfigItem{
				"topicNameProduce": {Value: "OrderReceived", DefaultValue: "OrderReceived", Alias: "broker_topic_name", Description: "broker's topic name used by the service"},
				"eventNameProduce": {Value: "OrderReceivedEvent", DefaultValue: "OrderReceivedEvent", Alias: "broker_event_name", Description: "broker's event name used by the service"},
				"topicNameDLQ":     {Value: "DeadLetterQueue", DefaultValue: "DeadLetterQueue", Alias: "broker_topic_name_dlq", Description: "DLQ's topic name used by the service to produce DLQ messages"},
				"eventNameDLQ":     {Value: "DeadLetterQueueEvent", DefaultValue: "DeadLetterQueueEvent", Alias: "broker_event_name_dlq", Description: "DLQ's event name used by the service to produce"},
			},
		},
		"zk": {
			Priority: 4,
			Items: map[string]config_helpers.ConfigItem{
				"host":           {Value: "127.0.0.1", DefaultValue: "127.0.0.1", Alias: "zk_host", Description: "zookeepeer host to connect"},
				"sessionTimeout": {Value: "1", DefaultValue: "1", Alias: "zk_session_timeout", Description: "zookeepeer session timeout"},
				"rootNode":       {Value: "/kafka_basics/Order", DefaultValue: "/kafka_basics/Order", Alias: "zk_root_node", Description: "zookeepeer root node for the service to read"},
			},
		},
		"api": {
			Priority: 4,
			Items: map[string]config_helpers.ConfigItem{
				"APISchemaFile": {Value: "api/swagger.yml", DefaultValue: "api/swagger.yml", Alias: "api_schema_file", Description: "schema file name to validate requests against"},
			},
		},
	},
}
