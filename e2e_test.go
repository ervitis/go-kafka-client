package gokafkaclient

import (
	"strconv"
	"sync"
	"testing"
)

func producerConfig() map[string]interface{} {
	return map[string]interface{}{
		"bootstrap.servers": "localhost:9092",
	}
}

func consumerConfig() map[string]interface{} {
	return map[string]interface{}{
		"group.id":           "test.simple-client",
		"auto.offset.reset":  "earliest",
		"session.timeout.ms": 10000,
		"bootstrap.servers":  "localhost:9092",
	}
}

func TestE2E(t *testing.T) {
	count := 0
	wg := sync.WaitGroup{}

	handler := func(msg []byte) {
		count++
		wg.Done()
	}

	errorHandler := func(msg []byte, err error) {}

	client := NewKafkaClient()

	consumer, err := client.SetConsumerConfig(consumerConfig()).BuildConsumer()
	if err != nil {
		t.Error(err)
	}

	producer, err := client.SetProducerConfig(producerConfig()).SetProducerTopicConfig("test-e2e", PartitionAny).BuildProducer()
	if err != nil {
		t.Error(err)
	}

	consumer.DeactivateValidator()
	producer.DeactivateValidator()

	wg.Add(5)
	go func() {
		consumer.Consume("test-e2e", handler, errorHandler)
	}()

	go func() {
		for i := 0; i < 5; i++ {
			_ = producer.Produce([]byte(`hello test ` + strconv.Itoa(i)))
		}
	}()

	wg.Wait()

	if count != 5 {
		t.Error("test e2e not worked")
	}
}
