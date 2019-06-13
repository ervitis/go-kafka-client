package go_kafka_client

import "testing"

func producerConfig() map[string]interface{} {
	return map[string]interface{}{
		"auto.commit": true,
		"poll.intervall.ms": 2000,
	}
}

func TestNewKafkaClient(t *testing.T) {
	if client := NewKafkaClient(); client == nil {
		t.Error("client is nil")
	}
}

func TestNewKafkaClient_NilProperties(t *testing.T) {
	client := NewKafkaClient()
	if client.cc == nil || client.pc == nil {
		t.Error("properties creating client nil")
	}
}

func TestKafkaClient_SetProducerConfig(t *testing.T) {
	client := NewKafkaClient()

	client.SetProducerConfig(producerConfig())

	if len(client.pc.config) == 0 {
		t.Error("config not set correctly")
	}
}

func TestKafkaClient_SetProducerConfig_Nil(t *testing.T) {
	client := NewKafkaClient()

	client.SetProducerConfig(nil)

	if client.pc.config != nil {
		t.Error("parameter config not set but the config has content")
	}
}

func TestKafkaClient_SetProducerConfig_Empty(t *testing.T) {
	client := NewKafkaClient()

	client.SetProducerConfig(map[string]interface{}{})

	if len(client.pc.config) != 0 {
		t.Error("parameter config empty but the config has content")
	}
}
