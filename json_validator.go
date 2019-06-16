package gokafkaclient

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"os"
)

type (
	JsonValidator struct{}
)

func (jv *JsonValidator) ValidateData(msg []byte, schema schema) (bool, error) {
	loader := gojsonschema.NewReferenceLoader(fmt.Sprintf("%s.%s", schema.Value, schema.Version))
	doc := gojsonschema.NewStringLoader(string(msg))

	v, err := gojsonschema.Validate(loader, doc)
	if err != nil {
		return false, err
	}

	var schErrors string
	for _, e := range v.Errors() {
		schErrors += fmt.Sprintf("%s; ", e.String())
	}

	return v.Valid(), fmt.Errorf(schErrors)
}

func (jv *JsonValidator) IsReachable(schema schema) bool {
	resp, err := http.Get(fmt.Sprintf("%s/%s.%s", os.Getenv("SCHEMA_HOST"), schema.Value, schema.Version))
	if err != nil {
		print(err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
