package apifootball

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

const (
	apiFootballDomain = "https://api-football-beta.p.rapidapi.com/"
)

// Client ...
type Client struct {
	domain     string
	options    *Options
	httpClient HttpClient
	// Endpoints	Endpoints
}

// Endpoint ..
type Endpoint interface {
	path() string
	params() string
}

// Query ...
type Query interface {
	verb() string
	path() string
}

var statusTimeout = 499

// NewClient ...
func NewClient(o *Options) *Client {
	return &Client{apiFootballDomain, o, &http.Client{}}
}

func request(c Client, q Query, data interface{}) (*http.Response, error) {
	resp, err := requestData(c, q)
	// This panics when the request has no response
	defer resp.Body.Close()

	if err != nil {
		return resp, err
	}

	decodeErr := decode(resp, &data)

	if decodeErr != nil {
		return resp, decodeErr
	}

	return resp, nil

}

func decode(resp *http.Response, data interface{}) error {
	jsonBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	jsonErr := json.Unmarshal(jsonBytes, &data)

	if jsonErr != nil {
		return jsonErr
	}

	return nil

}

func requestData(c Client, q Query) (*http.Response, error) {
	req, err := http.NewRequest(c.verb(q), c.url(q), nil)
	req.Header = c.headers()

	if err != nil {
		return nil, err
	}

	httpClient := c.httpClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	responseIsError := resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest

	if responseIsError {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, ErrUnauthorized
		case statusTimeout:
			return nil, ErrTimeout
		case http.StatusInternalServerError:
			return nil, ErrInternalServerError
		default:
			return nil, fmt.Errorf("unknown error: %s", resp.Status)
		}
	}

	return resp, nil
}

func (c Client) url(q Query) string {
	u, err := url.Parse(c.domain)
	if err != nil {
		log.Fatal(err)
	}

	u.Path = q.path()

	return u.String()
}

func (c Client) verb(q Query) string {
	return q.verb()
}

func (c Client) headers() http.Header {
	headers := http.Header{}
	headers.Set("x-rapidapi-key", c.options.key)

	return headers

}

// Fixtures ...
func (c Client) Fixtures(params map[string]interface{}) FixturesQuery {
	var result FixturesParameters
	err := mapstructure.Decode(params, &result)
	if err != nil {
		panic(err)
	}

	f := FixturesQuery{
		api:    c,
		Params: result,
	}

	return f
}
