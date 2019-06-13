package go_kafka_client

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func NewKafkaClient() *KafkaClient {
	return &KafkaClient{
		cc: &consumerClient{},
		pc: &producerClient{},
	}
}

func (kc *KafkaClient) SetProducerConfig(cfg map[string]interface{}) *KafkaClient {
	if cfg == nil {
		return kc
	}

	if kc.pc.config == nil {
		kc.pc.config = make(kafka.ConfigMap)
	}

	for k, v := range cfg {
		kc.pc.config[k] = v
	}

	return kc
}

func (kc *KafkaClient) SetConsumerConfig(cfg map[string]interface{}) *KafkaClient {
	if cfg == nil {
		return kc
	}

	if kc.cc.config == nil {
		kc.cc.config = make(kafka.ConfigMap)
	}

	for k, v := range cfg {
		kc.cc.config[k] = v
	}

	return kc
}

func (kc *KafkaClient) SetTimeoutPolling(polling int) *KafkaClient {
	if polling < defaultTimeout {
		kc.cc.pollTimeoutSeconds = defaultTimeout
	} else {
		kc.cc.pollTimeoutSeconds = polling
	}

	return kc
}

func (kc *KafkaClient) SetProducerTopicConfig(topicName string, partitionType int32) *KafkaClient {
	if kc.pc.t == nil {
		kc.pc.t = &kafka.TopicPartition{Topic: &topicName, Partition: partitionType}
	}

	return kc
}

func (kc *KafkaClient) BuildProducer() (*producerClient, error) {
	if kc.pc == nil || kc.pc.config == nil {
		return nil, fmt.Errorf(noConfigError, "producer")
	}

	if kc.pc.t == nil {
		return nil, fmt.Errorf(noConfigError, "producer for partition configuration")
	}

	if kc.pc.t.Topic == nil || *kc.pc.t.Topic == "" {
		return nil, fmt.Errorf(topicError, "producer", "topic name is empty")
	}

	var err error

	if kc.pc.p, err = kafka.NewProducer(&kc.pc.config); err != nil {
		return nil, err
	}

	return kc.pc, nil
}

func (kc *KafkaClient) BuildConsumer() (*consumerClient, error) {
	if kc.cc == nil || kc.cc.config == nil {
		return nil, fmt.Errorf(noConfigError, "consumer")
	}

	if kc.cc.pollTimeoutSeconds == 0 {
		kc.cc.pollTimeoutSeconds = defaultTimeout
	}

	var err error

	if kc.cc.c, err = kafka.NewConsumer(&kc.cc.config); err != nil {
		return nil, err
	}

	return kc.cc, nil
}
