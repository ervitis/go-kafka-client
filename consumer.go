package go_kafka_client

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
	"os/signal"
	"syscall"
)

func (c *consumerClient) ActivateValidator() *consumerClient {
	c.validateOnConsume = true
	return c
}

func (c *consumerClient) DeactivateValidator() *consumerClient {
	c.validateOnConsume = false
	return c
}

func (c *consumerClient) SetSchema(topic, schemaName, version string) *consumerClient {
	c.schemas[topic] = schema{Value: schemaName, Version: version}
	return c
}

func (c *consumerClient) dataIsValidFromSchema(ev []byte, sch schema) bool {
	if !c.validator.IsReachable(sch) {
		return false
	}

	if isValid, err := c.validator.ValidateData(ev, sch); err != nil {
		print(err)
		return false
	} else {
		return isValid
	}
}

func (c *consumerClient) Consume(topic string) error {
	if topic == "" {
		return fmt.Errorf(topicError, "consumer", "no topics to subscribe")
	}

	if err := c.c.Subscribe(topic, nil); err != nil {
		return err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	done := false
	var err error

	for !done {
		select {
		case sig := <-sigchan:
			err = fmt.Errorf(signalError, sig)
			done = true
		default:
			ev := c.c.Poll(c.pollTimeout)

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if c.validateOnConsume {
					if !c.dataIsValidFromSchema(e.Value, c.schemas[topic]) {
						return fmt.Errorf(schemaNotValidError, c.schemas[topic].Value, c.schemas[topic].Version, string(e.Value))
					}
				}
			case *kafka.Error:
				if e.IsFatal() {
					err = e
					done = true
				}
			}
		}
	}

	return err
}
