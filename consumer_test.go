package gokafkaclient

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"testing"
	"time"
)

type (
	mockConsumer struct{}
)

func (m *mockConsumer) receive(t time.Duration) (*kafka.Message, error) {
	return &kafka.Message{Value: []byte(`hello`)}, nil
}

func TestConsumerClient_ActivateValidator(t *testing.T) {
	c := consumerClient{}

	c.ActivateValidator()

	if !c.validateOnConsume {
		t.Error("validation on consume not set")
	}
}

func TestConsumerClient_DeactivateValidator(t *testing.T) {
	c := consumerClient{}

	c.ActivateValidator().DeactivateValidator()

	if c.validateOnConsume {
		t.Error("validation on consume set")
	}
}

func TestConsumerClient_SetSchema(t *testing.T) {
	c := consumerClient{}

	c.SetSchema("test", "test", "1.0.0")

	if v, ok := c.schemas["test"]; !ok {
		t.Error("no schema set")
	} else {
		if v.Value != "test" || v.Version != "1.0.0" {
			t.Error("schema values are not correct")
		}
	}
}

func TestConsumerClient_DataIsValidFromSchema(t *testing.T) {
	c := consumerClient{
		validator: &mockValidator{},
	}

	if !c.dataIsValidFromSchema([]byte(``), schema{Value: "valid", Version: "1.0.0"}) {
		t.Error("not valid when it is the schema")
	}
}

func TestConsumerClient_DataIsValidFromSchema_NotRechable(t *testing.T) {
	c := consumerClient{
		validator: &mockValidator{},
	}

	if c.dataIsValidFromSchema([]byte(``), schema{Value: "notreachable", Version: "1.0.0"}) {
		t.Error("not valid when it is not reachable")
	}
}

func TestConsumerClient_DataIsValidFromSchema_Error(t *testing.T) {
	c := consumerClient{
		validator: &mockValidator{},
	}

	if c.dataIsValidFromSchema([]byte(``), schema{Value: "error", Version: "1.0.0"}) {
		t.Error("on error it returned valid schema")
	}
}

func TestConsumerClient_DataIsValidFromSchema_NotValid(t *testing.T) {
	c := consumerClient{
		validator: &mockValidator{},
	}

	if c.dataIsValidFromSchema([]byte(``), schema{Value: "notvalid", Version: "1.0.0"}) {
		t.Error("schema is valid but it is not")
	}
}

func TestConsumerClient_filterEvent(t *testing.T) {
	c := consumerClient{}

	count := 0

	handler := func(msg []byte) {
		count++
	}

	conditions := [][]ConsumerConditions{
		{
			{Value: "1.0.0", Key: "VERSION"},
			{Value: "create-user", Key: "EventName"},
		},
		{
			{Value: "1.0.0", Key: "VERSION"},
		},
		{
			{Value: "2.0.0", Key: "VERSION"},
			{Value: "create-user", Key: "EventName"},
		},
		{
			{Value: "1.0.0", Key: "VERSIO"},
		},
	}

	c.filterEvent([]byte(`{"VERSION": "1.0.0", "EventName": "create-user", "Date": "2018-12-04"}`), handler, conditions[0])
	if count != 1 {
		t.Error("count != 1")
	}

	c.filterEvent([]byte(`{"VERSION": "1.0.0", "EventName": "create-user", "Date": "2018-12-04"}`), handler, conditions[1])
	if count != 2 {
		t.Error("count != 2")
	}

	c.filterEvent([]byte(`{"VERSION": "1.0.0", "EventName": "update-user", "Date": "2018-12-04"}`), handler, conditions[2])
	if count != 2 {
		t.Error("count not filter message!")
	}

	c.filterEvent([]byte(`{"VERSION": "2.0.0", "EventName": "create-user", "Date": "2018-12-04"}`), handler, conditions[2])
	if count != 3 {
		t.Error("count != 3")
	}

	c.filterEvent([]byte(`{"VERSION": "1.0.0", "EventName": "update-user", "Date": "2018-12-04"}`), handler, nil)
	if count != 4 {
		t.Error("count != 4")
	}

	c.filterEvent([]byte(`{"VERSION": "1.0.0", "EventName": "update-user", "Date": "2018-12-04"}`), handler, []ConsumerConditions{})
	if count != 5 {
		t.Error("count != 5")
	}

	c.filterEvent([]byte(`{"VERSION": "1.0.0", "EventName": "update-user", "Date": "2018-12-04"}`), handler, conditions[3])
	if count != 5 {
		t.Error("count != 5")
	}
}

func TestConsumerClient_Consume_Errors(t *testing.T) {
	c := consumerClient{
		validator: &mockValidator{},
		c:         &mockConsumer{},
		kc:        &kafka.Consumer{},
	}

	hasError := false

	handler := func(msg []byte) {}

	errorHandler := func(msg []byte, err error) {
		hasError = true
	}

	c.Consume("", handler, errorHandler)

	if !hasError {
		t.Error("there should be an error because topic was empty on consumer")
	}
}
