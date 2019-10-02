package gokafkaclient

import "github.com/confluentinc/confluent-kafka-go/kafka"

const (
	noConfigError = "no config set in %s"
	topicError    = "topic error in %s: %s"
	signalError   = "caught signal %s"

	schemaNotValidError = "schema %s with version %s is not valid with data %s"

	defaultTimeout = -1

	/**
	Type of partition, it will be handled by kafka library
	 */
	PartitionAny = kafka.PartitionAny
)
