package kafka_helpers

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (c *KafkaProducer) SetConfig(configMap *kafka.ConfigMap, produceTopic string) {
	c.Lock()
	c.configMap = configMap
	c.topic = produceTopic
	defer c.Unlock()
}

func (c *KafkaProducer) FilterConfig(configMap *kafka.ConfigMap) *kafka.ConfigMap {
	delete(*configMap, "group.id")
	return configMap
}

func (c *KafkaProducer) Create() error {
	var err error
	c.Lock()
	c.kafkaProducer, err = kafka.NewProducer(c.configMap)
	defer c.Unlock()
	return err
}

func (c *KafkaProducer) Close() {
	if c.kafkaProducer == nil {
		log.Println("KafkaProducer.Close(): Attempt of closing nil c.KafkaProducer")
		return
	}

	c.Lock()
	c.kafkaProducer.Close()
	defer c.Unlock()
}

func (c *KafkaProducer) Run() error {
	return nil
}

func (c *KafkaProducer) PublishEvent(event *BaseEvent, produceTimeoutMs int) error {
	log.Println("KafkaProducer.PublishEvent() is called")

	var value []byte
	var err error

	// TODO implement protobuf
	if value, err = json.Marshal(event); err != nil {
		return err
	}

	// Produce the message asynchronously
	if err = c.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &c.topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, nil); err != nil {
		log.Println("KafkaProducer.PublishEvent(): message could not be enqueued. Exiting...")
		return err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range c.kafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed producing message to Kafka: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			case kafka.Error:
				if ev.IsFatal() && ev.Code() == kafka.ErrLeaderNotAvailable {
					log.Printf("Kafka error: connection refused: %s\n", ev.String())
				} else {
					log.Printf("Kafka error: error producing message: %s\n", ev.String())
				}
			default:
				log.Printf("Ignored Kafka event: %v\n", e)
			}
		}
	}()

	// Wait for message deliveries before shutting down
	if outstandingNum := c.kafkaProducer.Flush(produceTimeoutMs); outstandingNum > 0 {
		log.Printf("Expected empty queue after Flush(), still has %d\n", outstandingNum)
	}

	return nil
}

func (c *KafkaProducer) GetTopic() string {
	return c.topic
}

func (c *KafkaProducer) Reload(configMap *kafka.ConfigMap, kafkaMessageHandlerFunc func(*kafka.Message)) {
	c.Close()
	c.SetConfig(configMap, c.GetTopic())
	if err := c.Create(); err != nil {
		log.Printf("KafkaProducer.Reload(): can't Create() KafkaProducer instance: %v\n", err.Error())
		return
	}
}
