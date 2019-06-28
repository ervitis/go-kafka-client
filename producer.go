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

	var hs []kafka.Header
	for _, header := range headers {
		h := kafka.Header{Key: header.Key, Value: []byte(header.Value)}
		hs = append(hs, h)
	}

	event, err := p.p.send(ev, hs...)
	if err != nil {
		return err
	}

	msg := event.(*kafka.Message)
	if err := msg.TopicPartition.Error; err != nil {
		return err
	}
	return nil
}

func (p *producerClient) send(ev []byte, headers ...kafka.Header) (kafka.Event, error) {
	devent := make(chan kafka.Event)

	defer close(devent)

	if err := p.kp.Produce(&kafka.Message{Headers: headers, TopicPartition: *p.t, Value: ev}, devent); err != nil {
		return nil, err
	}

	event := <-devent
	return event, nil
}
