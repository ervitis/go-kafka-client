package gokafkaclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/**
Activate the validator schema
 */
func (c *consumerClient) ActivateValidator() *consumerClient {
	c.validateOnConsume = true
	return c
}

/**
Desactivate the validator schema
 */
func (c *consumerClient) DeactivateValidator() *consumerClient {
	c.validateOnConsume = false
	return c
}

/**
ValidateMessageWithSchema if you want to validate the schema with the event data in []byte
 */
func (c *consumerClient) ValidateMessageWithSchema(data []byte, schema schema) bool {
	return c.dataIsValidFromSchema(data, schema)
}

/**
Set a schema to validate for the topic
 */
func (c *consumerClient) SetSchema(topic, schemaName, version string) *consumerClient {
	if c.schemas == nil {
		c.schemas = make(map[string]schema)
	}
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

/**
Subscribe to a topic parsing its handler and error handler and conditions
 */
func (c *consumerClient) Subscribe(topic string, handler ConsumerHandler, errHandler ConsumerErrorHandler, conditions ...ConsumerConditions) error {
	c.handler = handler
	c.errHandler = errHandler
	c.conditions = conditions

	if strings.TrimSpace(topic) == "" {
		return errors.New("cannot subscribe to empty topic")
	}
	c.topic = topic

	if err := c.kc.Subscribe(topic, nil); err != nil {
		return err
	}
	c.consume()
	return nil
}

/**
Consumer for the messages of the topic. When a message is read it will be filtered by the conditions and then the
handler will be called
 */
func (c *consumerClient) consume() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			err := fmt.Errorf(signalError, sig)
			c.errHandler(nil, err)
		default:
			msg, err := c.readMessage(time.Duration(c.pollTimeoutSeconds) * time.Second)
			if err != nil {
				switch erv := err.(type) {
				case *kafka.Error:
					if erv.IsFatal() {
						c.errHandler(nil, err)
					}
				default:
					c.errHandler(nil, err)
					return
				}
			}

			if c.validateOnConsume {
				if !c.dataIsValidFromSchema(msg.Value, c.schemas[c.topic]) {
					err := fmt.Errorf(schemaNotValidError, c.schemas[c.topic].Value, c.schemas[c.topic].Version, string(msg.Value))
					c.errHandler(msg.Value, err)
					return
				}
			}

			c.filterEvent(msg.Value, c.handler, c.conditions)
		}
	}
}

func (c *consumerClient) readMessage(t time.Duration) (*kafka.Message, error) {
	return c.kc.ReadMessage(t)
}

func (c *consumerClient) filterEvent(msg []byte, handler ConsumerHandler, conditions []ConsumerConditions) {
	var data map[string]interface{}
	_ = json.Unmarshal(msg, &data)

	n := len(conditions)
	count := 0

	if n == 0 {
		handler(msg)
	} else {
		for _, condition := range conditions {
			if _, ok := data[condition.Key]; !ok {
				continue
			}

			if data[condition.Key] == condition.Value {
				count++
			}
		}

		if n == count {
			handler(msg)
		}
	}
}
