package apifootball

import "net/http"

// Options is
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	key        string
	httpClient HttpClient
}
