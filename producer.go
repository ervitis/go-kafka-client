package gokafkaclient

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

/**
Activate the validator schema
*/
func (p *producerClient) ActivateValidator() *producerClient {
	p.validateOnProduce = true
	return p
}

/**
Deactivate the validator schema
*/
func (p *producerClient) DeactivateValidator() *producerClient {
	p.validateOnProduce = false
	return p
}

/**
Set a schema to validate for the topic
*/
func (p *producerClient) SetSchema(topic, schemaName, version string) *producerClient {
	if p.schemas == nil {
		p.schemas = make(map[string]schema)
	}
	p.schemas[topic] = schema{Value: schemaName, Version: version}
	return p
}

func (p *producerClient) dataIsValidFromSchema(ev []byte, schema schema) bool {
	if !p.validator.IsReachable(schema) {
		return false
	}

	if isValid, err := p.validator.ValidateData(ev, schema); err != nil {
		print(err)
	} else {
		return isValid
	}

	return false
}

/**
Produce an event with some optional headers
 */
func (p *producerClient) Produce(ev []byte, headers ...ProducerHeader) error {
	if p.t == nil {
		return fmt.Errorf(topicError, "producer", "topic config is empty")
	}

	devent := make(chan kafka.Event)

	var hs []kafka.Header
	for _, header := range headers {
		h := kafka.Header{Key: header.Key, Value: []byte(header.Value)}
		hs = append(hs, h)
	}

	defer close(devent)

	if err := p.p.Produce(&kafka.Message{Headers: hs, TopicPartition: *p.t, Value: ev}, devent); err != nil {
		return err
	}

	event := <-devent

	msg := event.(*kafka.Message)
	if err := msg.TopicPartition.Error; err != nil {
		return err
	}
	return nil
}
