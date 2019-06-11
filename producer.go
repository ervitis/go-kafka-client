package go_kafka_client

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func (p *producerClient) ActivateValidator() *producerClient {
	p.validateOnProduce = true
	return p
}

func (p *producerClient) DeactivateValidator() *producerClient {
	p.validateOnProduce = false
	return p
}

func (p *producerClient) SetSchema(topic, schemaName, version string) *producerClient {
	p.schemas[topic] = schema{Value: schemaName, Version: version}
	return p
}

func (p *producerClient) dataIsValidFromSchema(ev []byte, schema schema) bool {
	if !p.validator.IsReachable(schema) {
		return false
	}

	if isValid, err := p.validator.ValidateData(ev, schema); err != nil {
		panic(err)
	} else {
		return isValid
	}

	return false
}

func (p *producerClient) Produce(ev []byte, headers ...kafka.Header) error {
	if p.t == nil {
		return fmt.Errorf(topicError, "producer", "topic config is empty")
	}

	devent := make(chan kafka.Event)

	defer close(devent)

	if err := p.p.Produce(&kafka.Message{Headers: headers, TopicPartition: *p.t, Value: ev}, devent); err != nil {
		return err
	}

	event := <-devent

	msg := event.(*kafka.Message)
	if err := msg.TopicPartition.Error; err != nil {
		return err
	}
	return nil
}
