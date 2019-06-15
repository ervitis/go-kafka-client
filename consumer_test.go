package go_kafka_client

import (
	"testing"
)

type mockValidator struct {}

func (m *mockValidator) ValidateData(msg []byte, schema schema) (bool, error) {
	if schema.Value == "notvalid" {
		return false, nil
	}
	return true, nil
}

func (m *mockValidator) IsReachable(schema schema) bool {
	if schema.Value == "notreachable" {
		return false
	}
	return true
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
}