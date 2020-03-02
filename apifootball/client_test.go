package apifootball

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDecodeBlank(t *testing.T) {
	type Test struct {
		Field string `json:"field"`
	}

	data := Test{}
	err := decode(createFakeResponse(`{}`, 200), &data)
	expect(t, err, nil)
	expect(t, data.Field, "")
}

func TestDecodeValid(t *testing.T) {
	type Test struct {
		Field string `json:"field"`
	}

	data := Test{}
	err := decode(createFakeResponse(`{"field": "test"}`, 200), &data)
	expect(t, err, nil)
	expect(t, data.Field, "test")
}
func TestDecodeErr(t *testing.T) {
	type Test struct {
		Field string `json:"field"`
	}

	data := Test{}
	err := decode(createFakeResponse(``, 200), &data)
	expect(t, err.Error(), "unexpected end of JSON input")
}

type ATestQuery struct {
	verbVal string
	pathVal string
}

func (tq ATestQuery) verb() string { return tq.verbVal }
func (tq ATestQuery) path() string { return tq.pathVal }

func TestRequestData(t *testing.T) {
	query := ATestQuery{
		verbVal: "GET",
		pathVal: "/path/to/resource",
	}

	server, api := startMockServer(t, []RequestHandler{RequestHandler{
		PathAndQuery:  "/path/to/resource",
		ContentToSend: `{"test": "test"}`}})
	defer server.Close()

	response, _ := requestData(*api, query)

	expect(t, response.StatusCode, 200)
	jsonBytes, _ := ioutil.ReadAll(response.Body)
	json := string(jsonBytes)

	expect(t, json, "{\"test\": \"test\"}\n")

}

func TestRequest(t *testing.T) {
	type Test struct {
		Field string `json:"field"`
	}

	query := ATestQuery{
		verbVal: "GET",
		pathVal: "path/to/resource",
	}

	server, api := startMockServer(t, []RequestHandler{RequestHandler{
		PathAndQuery:  "/path/to/resource",
		ContentToSend: `{"field": "test"}`}})
	defer server.Close()

	data := Test{}
	response, _ := request(*api, query, &data)

	expect(t, data.Field, "test")
	expect(t, response.StatusCode, 200)
}

func TestClientUrl(t *testing.T) {
	query := ATestQuery{
		verbVal: "GET",
		pathVal: "path/to/resource",
	}

	client := NewClient(&Options{})
	client.domain = "http://some.domain"

	expect(t, client.url(query), "http://some.domain/path/to/resource")
}

func TestClientVerb(t *testing.T) {
	query := ATestQuery{
		verbVal: "VERB",
		pathVal: "path/to/resource",
	}

	client := NewClient(&Options{})

	expect(t, client.verb(query), "VERB")
}

func TestNewClient(t *testing.T) {
	client := NewClient(&Options{})

	expect(t, client.httpClient, &http.Client{})
}
