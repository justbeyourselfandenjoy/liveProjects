package kafka_helpers

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaGeneric struct {
	sync.Mutex
	configMap *kafka.ConfigMap
	topic     string
}

type KafkaProducer struct {
	KafkaGeneric
	kafkaProducer *kafka.Producer
}

type KafkaConsumer struct {
	KafkaGeneric
	kafkaConsumer  *kafka.Consumer
	kafkaMsgChan   chan *kafka.Message
	reloadInfoChan chan bool
}

type KafkaInstance interface {
	SetConfig(*kafka.ConfigMap, string)
	FilterConfig(*kafka.ConfigMap) *kafka.ConfigMap
	Create() error
	Close()
	Run() error
	GetTopic() string
	Reload(*kafka.ConfigMap, func(*kafka.Message))
}
