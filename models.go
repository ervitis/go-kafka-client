package gokafkaclient

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"time"
)

type (
	kafkaProducer interface {
		send(ev []byte, headers ...kafka.Header) (kafka.Event, error)
	}

	kafkaConsumer interface {
		receive(t time.Duration) (*kafka.Message, error)
	}

	producerClient struct {
		config            kafka.ConfigMap
		p                 kafkaProducer
		kp                *kafka.Producer
		t                 *kafka.TopicPartition
		validateOnProduce bool
		validator         Validator
		schemas           map[string]schema
	}

	consumerClient struct {
		config             kafka.ConfigMap
		c                  kafkaConsumer
		kc                 *kafka.Consumer
		pollTimeoutSeconds int
		validateOnConsume  bool
		validator          Validator
		schemas            map[string]schema
	}

	/**
	Kafka client
	*/
	KafkaClient struct {
		pc *producerClient
		cc *consumerClient
	}

	schema struct {
		Version string
		Value   string
	}

	/**
	Consumer filter conditions
	*/
	ConsumerConditions struct {
		Key   string
		Value string
	}

	/**
	Consumer handler type of function
	*/
	ConsumerHandler func(msg []byte)
	/**
	Consumer error handler
	*/
	ConsumerErrorHandler func(msg []byte, err error)

	/**
	Producer header
	*/
	ProducerHeader struct {
		Key   string
		Value string
	}
)
