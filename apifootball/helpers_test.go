package apifootball

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expect(t *testing.T, value interface{}, expected interface{}) {
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected %v  - Got %v (%T)", expected, value, value)
	}
}

func createFakeResponse(body string, statusCode int) *http.Response {
	return &http.Response{StatusCode: statusCode,
		Body: nopCloser{bytes.NewReader([]byte(body))}}
}

type RequestHandler struct {
	PathAndQuery string
	Method       string

	EstimatedContent string
	EstimatedHeaders map[string]string

	HeadersToSend    map[string]string
	ContentToSend    string
	StatusCodeToSend int
}

func startMockServer(t *testing.T, handlers []RequestHandler) (*httptest.Server, *Client) {
	options := Options{key: "Sup3r-Secret"}
	api := NewClient(&options)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			if handler.Method == "" {
				handler.Method = http.MethodGet
			}
			if handler.StatusCodeToSend == 0 {
				handler.StatusCodeToSend = http.StatusOK
			}
			if options.key != r.Header.Get("x-rapidapi-key") {
				panic("There is no api key set")
			}
			if handler.Method == r.Method && handler.PathAndQuery == r.URL.String() {

				if handler.EstimatedContent != "" {
					expect(t, readText(t, r.Body), handler.EstimatedContent)
				}
				if handler.EstimatedHeaders != nil {
					for key, value := range handler.EstimatedHeaders {
						expect(t, r.Header.Get(key), value)
					}
				}
				header := w.Header()
				if handler.HeadersToSend != nil {
					for key, value := range handler.HeadersToSend {
						header.Set(key, value)
					}
				}

				if handler.ContentToSend != "" && header.Get("Content-Type") == "" {
					header.Set("Content-Type", "application/json")
				}

				w.WriteHeader(handler.StatusCodeToSend)
				if handler.ContentToSend != "" {
					fmt.Fprintln(w, handler.ContentToSend)
				}
				return
			}
		}
		t.Logf("Unhandled request %s %s", r.Method, r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	}))

	api.domain = mockServer.URL

	return mockServer, api
}

func readText(t *testing.T, r io.Reader) string {
	text, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error("Error on reading content")
		return ""
	}
	return string(text)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
