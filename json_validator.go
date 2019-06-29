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
	loader := gojsonschema.NewReferenceLoader(fmt.Sprintf("%s.%s.json", schema.Value, schema.Version))
	doc := gojsonschema.NewStringLoader(string(msg))

	v, err := gojsonschema.Validate(loader, doc)
	if err != nil {
		return false, err
	}

	var schError error
	var s string
	for _, e := range v.Errors() {
		s += fmt.Sprintf("%s; ", e.String())
	}

	if s != "" {
		schError = fmt.Errorf(s)
	}

	return v.Valid(), schError
}

func (jv *JsonValidator) IsReachable(schema schema) bool {
	resp, err := http.Get(fmt.Sprintf("%s/%s.%s.json", os.Getenv("SCHEMA_HOST"), schema.Value, schema.Version))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return resp.StatusCode == http.StatusOK
}
