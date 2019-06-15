package go_kafka_client

import "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

const (
	noConfigError = "no config set in %s"
	topicError    = "topic error in %s: %s"
	signalError   = "caught signal %s"

	schemaNotValidError = "schema %s with version %s is not valid with data %s"

	defaultTimeout = -1

	PartitionAny = kafka.PartitionAny
)
