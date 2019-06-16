package go_kafka_client

import "errors"

type mockValidator struct{}

func (m *mockValidator) ValidateData(msg []byte, schema schema) (bool, error) {
	if schema.Value == "notvalid" {
		return false, nil
	} else if schema.Value == "error" {
		return false, errors.New("ups")
	}
	return true, nil
}

func (m *mockValidator) IsReachable(schema schema) bool {
	return schema.Value != "notreachable"
}
