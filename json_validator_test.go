package gokafkaclient

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var msg = []byte(`{"name": "test", "version": "1.0.0"}`)

func TestJsonValidator_ValidateData(t *testing.T) {
	schema := schema{Value: "file://./test_data/goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	ok, err := v.ValidateData(msg, schema)
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("the schema should be valid")
	}
}

func TestJsonValidator_ValidateDataErrorsSchema(t *testing.T) {
	schema := schema{Value: "file://./test_data/goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	ok, err := v.ValidateData([]byte(`{"version": "1.0.0"}`), schema)
	if err == nil {
		t.Error("error should not be empty")
	}

	if ok {
		t.Error("the schema should not be valid")
	}
}

func TestJsonValidator_ValidateDataErrorLoadingSchema(t *testing.T) {
	schema := schema{Value: "file:///test_data/goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	_, err := v.ValidateData(msg, schema)
	if err == nil {
		t.Error("error should not be empty")
	}
}

func startServer(port string, code int) *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadFile("./test_data/goodjson.1.0.0.json")

		w.WriteHeader(code)
		_, _ = w.Write(b)
	}))

	l, _ := net.Listen("tcp", ":" + port)
	server.Listener = l

	return server
}

func TestJsonValidator_IsReachable(t *testing.T) {
	server := startServer("8800", http.StatusOK)
	server.Start()
	defer server.Close()
	os.Setenv("SCHEMA_HOST", "http://localhost:8800")

	schema := schema{Value: "goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	if !v.IsReachable(schema) {
		t.Error("schema should be reachable")
	}
}

func TestJsonValidator_IsNotReachable(t *testing.T) {
	server := startServer("8800", http.StatusOK)
	server.Start()
	defer server.Close()
	os.Setenv("SCHEMA_HOST", "http://localhost:8806")

	schema := schema{Value: "goodjson", Version: "1.0.0"}

	v := &JsonValidator{}

	if v.IsReachable(schema) {
		t.Error("schema should not be reachable")
	}
}

func TestJsonValidator_IsNotReachable404(t *testing.T) {
	server := startServer("8800", http.StatusNotFound)
	server.Start()
	defer server.Close()
	os.Setenv("SCHEMA_HOST", "http://localhost:8800")
	schema := schema{Value: "goodson", Version: "1.1.0"}

	v := &JsonValidator{}

	if v.IsReachable(schema) {
		t.Error("schema should not be reachable and return 404")
	}
}
