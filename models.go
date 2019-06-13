package go_kafka_client

import "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

type (
	producerClient struct {
		config            kafka.ConfigMap
		p                 *kafka.Producer
		t                 *kafka.TopicPartition
		validateOnProduce bool
		validator         Validator
		schemas           map[string]schema
	}

	consumerClient struct {
		config            kafka.ConfigMap
		c                 *kafka.Consumer
		pollTimeout       int
		validateOnConsume bool
		validator         Validator
		schemas           map[string]schema
	}

	KafkaClient struct {
		pc *producerClient
		cc *consumerClient
	}

	schema struct {
		Version string
		Value   string
	}

	ConsumerConditions struct {
		Key   string
		Value string
	}

	ConsumerHandler func(msg []byte)
	ConsumerErrorHandler func(msg []byte, err error)
)
