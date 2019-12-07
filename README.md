# go-kafka-client

Kafka Golang library using confluent golang library and librdkafka

[![Build Status](https://travis-ci.org/ervitis/go-kafka-client.svg?branch=master)](https://travis-ci.org/ervitis/go-kafka-client) 
[![codecov](https://codecov.io/gh/ervitis/go-kafka-client/branch/master/graph/badge.svg)](https://codecov.io/gh/ervitis/go-kafka-client) 
[![Go Report Card](https://goreportcard.com/badge/github.com/ervitis/go-kafka-client)](https://goreportcard.com/report/github.com/ervitis/go-kafka-client) 
[![CircleCI](https://circleci.com/gh/ervitis/go-kafka-client/tree/master.svg?style=svg)](https://circleci.com/gh/ervitis/go-kafka-client/tree/master) 
[![Coverage Status](https://coveralls.io/repos/github/ervitis/go-kafka-client/badge.svg?branch=master)](https://coveralls.io/github/ervitis/go-kafka-client?branch=master)

## Getting Started

You can get the library using `go get`

```bash
go get github.com/ervitis/go-kafka-client
```

In this repository there is a `docker-compose.yml` repository where we can try the library using a docker image from `landoop`

For an example how to use this library, see inside the folder `examples`.

We can try it if you run the following commands:

```bash
go run examples/simple_consumer/main.go
```

Open another terminal and type:

```bash
go run examples/simple_producer/main.go
```

To stop the consumer send a kill signal to the process

### API

#### Kafka client

To build a client just call the function

```go
package main

import gkc "github.com/ervitis/go-kafka-client"

client := gkc.NewKafkaClient()
```

Created the client we can set the configuration for the kafka producer and consumer

```go
producer := client.SetProducerConfig(...)
                  .SetProducerTopicConfig("topic", gkc.PartitionAny)
                  .BuildProducer()

consumer := client.SetConsumerConfig(...).BuildConsumer()
```

#### Producer

When using the producer, we can set it schema to validate the message it will produce.

```go
producer.SetSchema("topic", "create-user", "1.0.0").ActivateValidator()
```

Then, produce an event

```go
err := producer.Produce([]byte(`hello`))
if err != nil {
	panic(err)
}
```

#### Consumer

Using the consumer function is simple as using the producer, it need the topic from it will  
listen the event, two handler functions and an array of condition to filter the events.

The consumer will be waiting until it can read a message, so we should use `go routines` to handle them.

```go
func handlerEvent(msg []byte) { fmt.Println("ey!") }

func handlerError(msg []byte, err error) { fmt.Println(err) }

conditions := []ConsumerConditions{ {Key: "EventType", Value: "create-user"} }

consumer.Subscribe("topic", handlerEvent, handlerError, conditions)
```

### Prerequisites

This library is using confluent kafka go client version 1.1.0 and librdkafka 1.3.0 version

```bash
git clone -b v1.0.1 https://github.com/edenhill/librdkafka

./configure
make
sudo make install
```

## Running the tests

```bash
go test -race -v ./...
```

## Built With

* [Golang](http://www.golang.org) - GoLang programming language

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Victor Martin** - *Initial work* - [ervitis](https://github.com/ervitis)

## License

This project is licensed under the Apache 2.0 - see the [LICENSE.md](LICENSE.md) file for details
