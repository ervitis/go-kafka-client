package gokafkaclient

type (
	Validator interface {
		ValidateData(msg []byte, schema schema) (bool, error)
		IsReachable(schema schema) bool
	}
)
