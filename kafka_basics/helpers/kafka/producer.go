package kafka_helpers

import "github.com/confluentinc/confluent-kafka-go/kafka"

func CreateKafkaProducerInstance(configMap *kafka.ConfigMap) (*kafka.Producer, error) {
	return kafka.NewProducer(configMap)
}
