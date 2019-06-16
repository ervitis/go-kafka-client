package main

import (
	"fmt"
	kcl "github.com/ervitis/go-kafka-client"
)

func producerConfig() map[string]interface{} {
	return map[string]interface{}{
		"bootstrap.servers": "localhost:9092",
	}
}

func main() {
	client := kcl.NewKafkaClient()

	topic := "testing.simple"

	p, err := client.SetProducerTopicConfig(topic, kcl.PartitionAny).SetProducerConfig(producerConfig()).BuildProducer()
	if err != nil {
		panic(err)
	}

	p.DeactivateValidator()

	for i := 0; i < 50; i++ {
		ev := []byte(fmt.Sprintf(`My message is inserted in order %d`, i))
		p.Produce(ev)
	}
}
