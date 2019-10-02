package gokafkaclient

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"testing"
)

type (
	mockProducer struct{}

	mockBadProducer struct{}

	mockTopicProducerError struct{}
)

func (m *mockProducer) send(kc *kafka.Producer, tp *kafka.TopicPartition, ev []byte, headers ...kafka.Header) (kafka.Event, error) {
	mevent := &kafka.Message{}
	return mevent, nil
}

func (m *mockBadProducer) send(kc *kafka.Producer, tp *kafka.TopicPartition, ev []byte, headers ...kafka.Header) (kafka.Event, error) {
	mevent := &kafka.Message{}
	return mevent, errors.New("ups")
}

func (m *mockTopicProducerError) send(kc *kafka.Producer, tp *kafka.TopicPartition, ev []byte, headers ...kafka.Header) (kafka.Event, error) {
	mevent := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Error: errors.New("ups")},
	}
	return mevent, nil
}

func TestProducerClient_ActivateValidator(t *testing.T) {
	c := producerClient{}

	c.ActivateValidator()

	if !c.validateOnProduce {
		t.Error("validation on consume not set")
	}
}

func TestProducerClient_DeactivateValidator(t *testing.T) {
	c := producerClient{}

	c.ActivateValidator().DeactivateValidator()

	if c.validateOnProduce {
		t.Error("validation on consume set")
	}
}

func TestProducerClient_SetSchema(t *testing.T) {
	c := producerClient{}

	c.SetSchema("test", "test", "1.0.0")

	if v, ok := c.schemas["test"]; !ok {
		t.Error("no schema set")
	} else {
		if v.Value != "test" || v.Version != "1.0.0" {
			t.Error("schema values are not correct")
		}
	}
}

func TestProducerClient_DataIsValidFromSchema(t *testing.T) {
	c := producerClient{
		validator: &mockValidator{},
	}

	if !c.dataIsValidFromSchema([]byte(``), schema{Value: "valid", Version: "1.0.0"}) {
		t.Error("not valid when it is the schema")
	}
}

func TestProducerClient_DataIsValidFromSchema_NotRechable(t *testing.T) {
	c := producerClient{
		validator: &mockValidator{},
	}

	if c.dataIsValidFromSchema([]byte(``), schema{Value: "notreachable", Version: "1.0.0"}) {
		t.Error("not valid when it is not reachable")
	}
}

func TestProducerClient_DataIsValidFromSchema_Error(t *testing.T) {
	c := producerClient{
		validator: &mockValidator{},
	}

	if c.dataIsValidFromSchema([]byte(``), schema{Value: "error", Version: "1.0.0"}) {
		t.Error("on error it returned valid schema")
	}
}

func TestProducerClient_DataIsValidFromSchema_NotValid(t *testing.T) {
	c := producerClient{
		validator: &mockValidator{},
	}

	if c.dataIsValidFromSchema([]byte(``), schema{Value: "notvalid", Version: "1.0.0"}) {
		t.Error("schema is valid but it is not")
	}
}

func TestProducerClient_Produce(t *testing.T) {
	topic := "test-producer"

	c := producerClient{
		validator: &mockValidator{},
		p:         &mockProducer{},
		t:         &kafka.TopicPartition{Topic: &topic, Partition: PartitionAny},
	}

	if err := c.Produce([]byte(`hello`)); err != nil {
		t.Error(err)
	}
}

func TestProducerClient_Produce_TopicError(t *testing.T) {
	c := producerClient{
		validator: &mockValidator{},
		p:         &mockProducer{},
		t:         nil,
	}

	if err := c.Produce([]byte(`hello`)); err == nil {
		t.Error("there should be an error if topic is nil")
	}
}

func TestProducerClient_Produce_WithHeaders(t *testing.T) {
	topic := "test-producer"

	c := producerClient{
		validator: &mockValidator{},
		p:         &mockProducer{},
		t:         &kafka.TopicPartition{Topic: &topic, Partition: PartitionAny},
	}

	headers := []ProducerHeader{
		{
			Key:   "header1",
			Value: "valueheader1",
		},
	}

	if err := c.Produce([]byte(`hello`), headers...); err != nil {
		t.Error("there was an error using headers")
	}
}

func TestProducerClient_Produce_ErrorOnTopicPartition(t *testing.T) {
	topic := "test-producer"

	c := producerClient{
		validator: &mockValidator{},
		p:         &mockTopicProducerError{},
		t:         &kafka.TopicPartition{Topic: &topic, Partition: PartitionAny},
	}

	headers := []ProducerHeader{
		{
			Key:   "header1",
			Value: "valueheader1",
		},
	}

	if err := c.Produce([]byte(`hello`), headers...); err == nil {
		t.Error("there should be an error producing")
	}
}

func TestProducerClient_Produce_ErrorProducer(t *testing.T) {
	topic := "test-producer"

	c := producerClient{
		validator: &mockValidator{},
		p:         &mockBadProducer{},
		t:         &kafka.TopicPartition{Topic: &topic, Partition: PartitionAny},
	}

	headers := []ProducerHeader{
		{
			Key:   "header1",
			Value: "valueheader1",
		},
	}

	if err := c.Produce([]byte(`hello`), headers...); err == nil {
		t.Error("there should be an error producing")
	}
}
