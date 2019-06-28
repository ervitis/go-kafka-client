package gokafkaclient

import "testing"

var msg = []byte(`{"name": "test", "version": "1.0.0"}`)

func TestJsonValidator_ValidateData(t *testing.T) {
	schema := schema{Value: "file://./test_data/goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	ok, err := v.ValidateData(msg, schema)
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("the schema should be valid")
	}
}

func TestJsonValidator_ValidateDataErrorsSchema(t *testing.T) {
	schema := schema{Value: "file://./test_data/goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	ok, err := v.ValidateData([]byte(`{"version": "1.0.0"}`), schema)
	if err == nil {
		t.Error("error should not be empty")
	}

	if ok {
		t.Error("the schema should not be valid")
	}
}

func TestJsonValidator_ValidateDataErrorLoadingSchema(t *testing.T) {
	schema := schema{Value: "file:///test_data/goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	_, err := v.ValidateData(msg, schema)
	if err == nil {
		t.Error("error should not be empty")
	}
}
