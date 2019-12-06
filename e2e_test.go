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

var count = 0
var wg sync.WaitGroup

var handler = func(msg []byte) {
	count++
	wg.Done()
}

var errorHandler = func(msg []byte, err error) {}

func TestE2E(t *testing.T) {

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

	for i := 0; i < 5; i++ {
		_ = producer.Produce([]byte(`hello test ` + strconv.Itoa(i)))
	}

	wg.Wait()

	wg.Add(5)
	go func() {
		_ = consumer.Subscribe("test-e2e", handler, errorHandler)
		consumer.Consume()
	}()

	wg.Wait()
	if count != 5 {
		t.Error("test e2e not worked")
	}
}
