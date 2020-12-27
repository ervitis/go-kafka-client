package tests

import (
	"context"
	"fmt"
	"github.com/cabify/aceptadora"
	"github.com/ervitis/go-kafka-client"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/suite"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
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

type (
	Config struct {
		Aceptadora  aceptadora.Config
		ImagePuller aceptadora.ImagePullerConfig
	}

	TestKafka struct {
		suite.Suite
		aceptadora *aceptadora.Aceptadora
		config     Config
	}
)

func (tk *TestKafka) SetupTest() {
	aceptadora.SetEnv(tk.T())

	tk.Require().NoError(envconfig.Process("ACCEPTANCE", &tk.config))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	imagePuller := aceptadora.NewImagePuller(tk.T(), tk.config.ImagePuller)
	tk.aceptadora = aceptadora.New(tk.T(), imagePuller, tk.config.Aceptadora)
	tk.aceptadora.PullImages(ctx)

	tk.aceptadora.Run(ctx, "kafka")
	tk.Require().Eventually(func() bool {
		ping, err := net.Dial("tcp", "localhost:9092")
		if err != nil {
			return false
		}
		_ = ping.Close()
		return true
	}, time.Minute, 3 * time.Second, "broker1 not started")
}

func (tk *TestKafka) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tk.aceptadora.StopAll(ctx)
}

var count = 0
var wg sync.WaitGroup

const N = 5

var handler = func(msg []byte) {
	fmt.Println("Before " + strconv.Itoa(count))
	count++
	fmt.Println("After " + strconv.Itoa(count))
	wg.Done()
}

var errorHandler = func(msg []byte, err error) {}

func (tk *TestKafka) TestE2E() {
	client := gokafkaclient.NewKafkaClient()

	consumer, err := client.SetConsumerConfig(consumerConfig()).BuildConsumer()
	if err != nil {
		tk.Error(err)
	}

	producer, err := client.SetProducerConfig(producerConfig()).SetProducerTopicConfig("test-e2e", gokafkaclient.PartitionAny).BuildProducer()
	if err != nil {
		tk.Error(err)
	}

	consumer.DeactivateValidator()
	producer.DeactivateValidator()

	ready := make(chan bool, 1)

	go func(ready chan bool) {
		if err = consumer.Subscribe("test-e2e", handler, errorHandler); err != nil {
			panic(err)
		}
		time.Sleep(15*time.Second)
		ready<-true
	}(ready)

	<-ready
	wg.Add(N)

	for i := 0; i < N; i++ {
		err = producer.Produce([]byte(`hello test ` + strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}
	}


	wg.Wait()
}

func TestAcceptance(t *testing.T) {
	suite.Run(t, new(TestKafka))
}
