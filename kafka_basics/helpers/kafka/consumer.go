package kafka_helpers

import (
	"errors"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (c *KafkaConsumer) SetConfig(configMap *kafka.ConfigMap, consumeTopic string) {
	c.Lock()
	c.configMap = configMap
	c.topic = consumeTopic
	defer c.Unlock()
}

func (c *KafkaConsumer) FilterConfig(configMap *kafka.ConfigMap) *kafka.ConfigMap {
	delete(*configMap, "request.required.acks")
	delete(*configMap, "message.timeout.ms")
	delete(*configMap, "go.delivery.report.fields")
	delete(*configMap, "go.delivery.reports")
	delete(*configMap, "go.batch.producer")
	delete(*configMap, "go.events.channel.size")
	delete(*configMap, "go.produce.channel.size")
	delete(*configMap, "go.logs.channel.enable")
	delete(*configMap, "go.logs.channel")
	return configMap
}

func (c *KafkaConsumer) Create() error {
	var err error
	c.Lock()
	if c.kafkaConsumer, err = kafka.NewConsumer(c.configMap); err == nil {
		c.createMsgChan()
		c.createReloadInfoChan()

		if err := c.kafkaConsumer.SubscribeTopics([]string{c.topic}, nil); err != nil {
			log.Println("KafkaConsumer.Create(): Failed to subscribe to the topic ", c.topic)
			defer c.Unlock()
			return err
		}
	}
	defer c.Unlock()
	return err
}

func (c *KafkaConsumer) Close() {
	if c.kafkaConsumer == nil {
		log.Println("KafkaConsumer.Close(): Attempt of closing nil c.kafkaConsumer")
		return
	}

	c.getReloadInfoChan() <- true
	c.Lock()
	if c.GetMsgChan() != nil {
		close(c.GetMsgChan())
	}

	if c.getReloadInfoChan() != nil {
		close(c.getReloadInfoChan())
	}

	c.kafkaConsumer.Unsubscribe()

	if err := c.kafkaConsumer.Close(); err != nil {
		log.Println("KafkaConsumer.Close(): Failed to close kafka instance")
	}
	defer c.Unlock()
}

func (c *KafkaConsumer) Run() error {
	var msg *kafka.Message
	var err error
	c.Lock()
	for !c.kafkaConsumer.IsClosed() {
		msg, err = c.kafkaConsumer.ReadMessage(time.Second)
		if err == nil {
			log.Printf("KafkaConsumer.Run(): Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			c.GetMsgChan() <- msg
		} else if err.(kafka.Error).Code() != kafka.ErrTimedOut {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			log.Printf("KafkaConsumer.Run(): Kafka consumer error: %v\n", err)
			defer c.Unlock()
			return err
		}

		select {
		case <-c.getReloadInfoChan():
			log.Println("KafkaConsumer.Run(): external reload signal is caught. Exiting...")
			defer c.Unlock()
			return nil
		default:
		}
	}
	defer c.Unlock()
	return errors.New("c.kafkaConsumer.IsClosed() = true, expect false")
}

func (c *KafkaConsumer) GetMsgChan() chan *kafka.Message {
	return c.kafkaMsgChan
}

func (c *KafkaConsumer) createMsgChan() {
	c.kafkaMsgChan = make(chan *kafka.Message, 1)
}

func (c *KafkaConsumer) getReloadInfoChan() chan bool {
	return c.reloadInfoChan
}

func (c *KafkaConsumer) createReloadInfoChan() {
	c.reloadInfoChan = make(chan bool)
}

func (c *KafkaConsumer) GetTopic() string {
	return c.topic
}

func (c *KafkaConsumer) Reload(configMap *kafka.ConfigMap, kafkaMessageHandlerFunc func(*kafka.Message)) {
	c.Close()
	c.SetConfig(configMap, c.GetTopic())
	if err := c.Create(); err != nil {
		log.Printf("KafkaConsumer.Reload(): can't Create() KafkaConsumer instance: %v\n", err.Error())
		return
	}

	go func() {
		if err := c.Run(); err != nil {
			<-c.getReloadInfoChan()
			log.Printf("KafkaConsumer.Reload(): can't Run() KafkaConsumer instance: %v\n", err.Error())
		}
	}()

	go func() {
		for {
			kafkaMessageHandlerFunc(<-c.GetMsgChan())
		}
	}()
}
