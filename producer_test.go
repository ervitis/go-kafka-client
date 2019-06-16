package go_kafka_client

import "testing"

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
