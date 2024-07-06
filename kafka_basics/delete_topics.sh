#!/bin/bash

KAFKA_BIN=/opt/kafka/bin
KAFKA_TOPICS="$KAFKA_BIN/kafka-topics.sh --bootstrap-server localhost:9092"
KAFKA_TOPIC_DELETE="$KAFKA_TOPICS --delete --topic"

#list topics before creation
echo "Existing topics:"
$KAFKA_TOPICS --list

#create topics
echo "Deleting the topics:"
$KAFKA_TOPIC_DELETE OrderReceived
$KAFKA_TOPIC_DELETE OrderConfirmed
$KAFKA_TOPIC_DELETE OrderPickedAndPacked
$KAFKA_TOPIC_DELETE Notification
$KAFKA_TOPIC_DELETE DeadLetterQueue

#list topics after creation
echo "All the topics after creation:"
$KAFKA_TOPICS --list
