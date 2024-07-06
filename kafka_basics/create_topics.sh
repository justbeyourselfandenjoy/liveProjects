#!/bin/bash

KAFKA_BIN=/opt/kafka/bin
KAFKA_TOPICS="$KAFKA_BIN/kafka-topics.sh --bootstrap-server localhost:9092"
KAFKA_TOPIC_CREATE="$KAFKA_TOPICS --create --topic"

#list topics before creation
echo "Existing topics:"
$KAFKA_TOPICS --list

#create topics
echo "Creating new topics:"
$KAFKA_TOPIC_CREATE OrderReceived
$KAFKA_TOPIC_CREATE OrderConfirmed
$KAFKA_TOPIC_CREATE OrderPickedAndPacked
$KAFKA_TOPIC_CREATE Notification
$KAFKA_TOPIC_CREATE DeadLetterQueue

#list topics after creation
echo "All the topics after creation:"
$KAFKA_TOPICS --list
