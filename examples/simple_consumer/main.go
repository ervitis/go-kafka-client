package main

import (
	"fmt"
	kcl "github.com/ervitis/go-kafka-client"
)

func consumerConfig() map[string]interface{} {
	return map[string]interface{}{
		"group.id":           "test.simple-client",
		"auto.offset.reset":  "earliest",
		"session.timeout.ms": 10000,
		"bootstrap.servers":  "localhost:9092",
	}
}

func handler(msg []byte) {
	fmt.Println(string(msg))
}

func errorHandler(msg []byte, err error) {
	fmt.Printf("%s with error %s", string(msg), err.Error())
}

func main() {
	client := kcl.NewKafkaClient()

	c, err := client.SetConsumerConfig(consumerConfig()).BuildConsumer()
	if err != nil {
		panic(err)
	}

	c.DeactivateValidator()

	topic := "testing.simple"

	c.Subscribe(topic, handler, errorHandler)
	c.Consume()
}
