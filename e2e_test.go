package gokafkaclient

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func producerConfig() map[string]interface{} {
	return map[string]interface{}{
		"bootstrap.servers": "localhost:9092",
		"acks":              -1,
	}
}

func consumerConfig() map[string]interface{} {
	return map[string]interface{}{
		"group.id":             "test.simple-client",
		"auto.offset.reset":    "earliest",
		"session.timeout.ms":   10000,
		"bootstrap.servers":    "localhost:9092",
	}
}

var count = 0
var wg sync.WaitGroup
var mtx sync.Mutex

const N = 5

var handler = func(msg []byte) {
	mtx.Lock()
	defer mtx.Unlock()

	if count < 5 {
		fmt.Println("Before " + strconv.Itoa(count))
		count++
		fmt.Println("After " + strconv.Itoa(count))
		wg.Done()
	}
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

	for i := 0; i < N; i++ {
		_ = producer.Produce([]byte(`hello test ` + strconv.Itoa(i)))
	}

	if err = consumer.Subscribe("test-e2e", handler, errorHandler); err != nil {
		panic(err)
	}
	go consumer.Consume()

	wg.Add(N)
	wg.Wait()
}
