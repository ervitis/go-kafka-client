package go_kafka_client

type (
	Validator interface {
		ValidateData(msg []byte, schema schema) (bool, error)
		IsReachable(schema schema) bool
	}
)
