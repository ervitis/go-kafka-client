package gokafkaclient

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"testing"
)

func exampleConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"log_level":               5,
		"socket.keepalive.enable": true,
	}
}

func consumerConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"group.id": "test",
	}
}

func badConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"ups": "wroooong",
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

	client.SetProducerConfig(exampleConfiguration())

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

func TestKafkaClient_SetConsumerConfig(t *testing.T) {
	client := NewKafkaClient()

	client.SetConsumerConfig(exampleConfiguration())

	if len(client.cc.config) == 0 {
		t.Error("config not set correctly")
	}
}

func TestKafkaClient_SetConsumerConfig_Nil(t *testing.T) {
	client := NewKafkaClient()

	client.SetConsumerConfig(nil)

	if client.cc.config != nil {
		t.Error("parameter config not set but the config has content")
	}
}

func TestKafkaClient_SetConsumerConfig_Empty(t *testing.T) {
	client := NewKafkaClient()

	client.SetConsumerConfig(map[string]interface{}{})

	if len(client.cc.config) != 0 {
		t.Error("parameter config empty but the config has content")
	}
}

func TestKafkaClient_SetTimeoutPolling(t *testing.T) {
	client := NewKafkaClient()

	client.SetTimeoutPolling(10)

	if client.cc.pollTimeoutSeconds != 10 {
		t.Error("timeout polling value distinct of 10")
	}
}

func TestKafkaClient_SetTimeoutPolling_NotValidValue(t *testing.T) {
	client := NewKafkaClient()

	client.SetTimeoutPolling(-2000)

	if client.cc.pollTimeoutSeconds != defaultTimeout {
		t.Error("timeout invalid value not set to default")
	}
}

func TestKafkaClient_SetProducerTopicConfig(t *testing.T) {
	client := NewKafkaClient().SetProducerTopicConfig("test", kafka.PartitionAny)

	if client.pc.t == nil {
		t.Error("producer client topic config is nil")
	}
}

func TestKafkaClient_BuildProducer(t *testing.T) {
	producer, err := NewKafkaClient().SetProducerConfig(exampleConfiguration()).SetProducerTopicConfig("test", kafka.PartitionAny).BuildProducer()
	if err != nil {
		t.Error(err)
	}

	if producer.kp == nil {
		t.Error("producer built nil")
	}
}

func TestKafkaClient_BuildProducer_NoConfigSet(t *testing.T) {
	_, err := NewKafkaClient().BuildProducer()

	if err == nil {
		t.Error("error not returned building with no config")
	}
}

func TestKafkaClient_BuildProducer_PartitionConfigNotSet(t *testing.T) {
	_, err := NewKafkaClient().SetProducerConfig(exampleConfiguration()).BuildProducer()

	if err == nil {
		t.Error("error not returned building with not partition config set")
	}
}

func TestKafkaClient_BuildProducer_EmptyTopic(t *testing.T) {
	_, err := NewKafkaClient().SetProducerConfig(exampleConfiguration()).SetProducerTopicConfig("", kafka.PartitionAny).BuildProducer()

	if err == nil {
		t.Error("error not returned building with not partition config set")
	}
}

func TestKafkaClient_BuildProducer_LibraryError(t *testing.T) {
	_, err := NewKafkaClient().SetProducerConfig(badConfiguration()).SetProducerTopicConfig("test", kafka.PartitionAny).BuildProducer()

	if err == nil {
		t.Error("error not returned building kafka library producer")
	}
}

func TestKafkaClient_BuildConsumer(t *testing.T) {
	consumer, err := NewKafkaClient().SetConsumerConfig(consumerConfiguration()).BuildConsumer()
	if err != nil {
		t.Error(err)
	}

	if consumer.kc == nil {
		t.Error("consumer built nil")
	}
}

func TestKafkaClient_BuildConsumer_NotSetPollTimeout(t *testing.T) {
	consumer, _ := NewKafkaClient().SetConsumerConfig(consumerConfiguration()).BuildConsumer()

	if consumer.pollTimeoutSeconds != defaultTimeout {
		t.Error("timeout not default")
	}
}

func TestKafkaClient_BuildConsumer_NotConfigurationSet(t *testing.T) {
	_, err := NewKafkaClient().BuildConsumer()
	if err == nil {
		t.Error("error not triggered if no config set in consumer build function")
	}
}

func TestKafkaClient_BuildConsumer_LibraryError(t *testing.T) {
	_, err := NewKafkaClient().SetConsumerConfig(badConfiguration()).BuildConsumer()

	if err == nil {
		t.Error("error not returned building kafka library consumer")
	}
}
