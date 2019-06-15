# go-kafka-client

Kafka Golang library using confluent golang library and librdkafka

## Getting Started

You can get the library using `go get`

```bash
go get github.com/ervitis/go-kafka-client
```

### Prerequisites

This library is using confluent kafka go client version 1.0.0 and librdkafka 1.0.1 version

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
