package go_kafka_client

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func (c *consumerClient) Consume(topic string, handler ConsumerHandler, errHandler ConsumerErrorHandler, conditions []ConsumerConditions) {
	if topic == "" {
		err := fmt.Errorf(topicError, "consumer", "no topics to subscribe")
		errHandler(nil, err)
		return
	}

	if err := c.c.Subscribe(topic, nil); err != nil {
		errHandler(nil, err)
		return
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for true {
		select {
		case sig := <-sigchan:
			err := fmt.Errorf(signalError, sig)
			errHandler(nil, err)
		default:
			msg, err := c.c.ReadMessage(time.Duration(c.pollTimeoutSeconds)*time.Second)
			if err != nil {
				errHandler(nil, err)
				return
			}

			if c.validateOnConsume {
				if !c.dataIsValidFromSchema(msg.Value, c.schemas[topic]) {
					err := fmt.Errorf(schemaNotValidError, c.schemas[topic].Value, c.schemas[topic].Version, string(msg.Value))
					errHandler(msg.Value, err)
					return
				}
			}

			c.filterEvent(msg.Value, handler, conditions)
		}
	}
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
