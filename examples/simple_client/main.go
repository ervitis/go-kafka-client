package main

import (
	"fmt"
	kcl "github.com/ervitis/go-kafka-client"
)

func main() {
	client := kcl.NewKafkaClient()

	c, err := client.BuildConsumer()
	if err != nil {
		fmt.Println(err)
	}

	print(c)

	print(client)
}
