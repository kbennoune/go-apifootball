package apifootball

import (
	"strconv"
	"time"
)

// FixturesQuery ...
type FixturesQuery struct {
	Params FixturesParameters
	api    Client
}

// FixturesEnvelope ...
type FixturesEnvelope struct {
	Get        string                 `json:"get"`
	Parameters map[string]interface{} `json:"parameters"`
	// errors	[]
	Results  int        `json:"results"`
	Response []Fixtures `json:"response"`
}

// FixturesParameters ...
type FixturesParameters struct {
	ID       int    `json:"id"`
	Live     string `json:"live"`
	Date     string `json:"date"`
	League   int    `json:"league"`
	Season   int    `json:"season"`
	Team     int    `json:"team"`
	Last     int    `json:"last"`
	Next     int    `json:"next"`
	From     string `json:"from"`
	To       string `json:"to"`
	Round    string `json:"round"`
	Timezone string `json:"timezone"`
}

// Fixtures ...
type Fixtures struct {
	Fixture Fixture `json:"fixture"`
	League  League  `json:"league"`
	Teams   Teams   `json:"teams"`
	Goals   Goals   `json:"goals"`
	Score   Score   `json:"score"`
}

// Fixture ...
type Fixture struct {
	ID       int    `json:"id"`
	Venue    string `json:"venue"`
	Referee  string `json:"referee"`
	Timezone string `json:"timezone"`
	// Date     time.Time `json:"date"`
	// Timestamp time.Time     `json:"timestamp"`
	Periods Periods       `json:"periods"`
	Status  FixtureStatus `json:"status"`
}

// Periods ...
type Periods struct {
	First  TimestampTime `json:"first"`
	Second TimestampTime `json:"second"`
}

// TimestampTime ...
type TimestampTime time.Time

// UnmarshalJSON ...
func (t *TimestampTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	secs, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return err
	}

	*t = TimestampTime(time.Unix(secs, 0))

	return nil
}

// FixtureStatus ...
type FixtureStatus struct {
	Long    string `json:"long"`
	Short   string `json:"short"`
	Elapsed int    `json:"elapsed"`
}

// League ...
type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int    `json:"season"`
	Round   string `json:"round"`
}

// Team ...
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// Teams ...
type Teams struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

// Goals ...
type Goals struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Score ...
type Score struct {
	Halftime  Goals `json:"halftime"`
	Fulltime  Goals `json:"fulltime"`
	Extratime Goals `json:"extratime"`
	Penalty   Goals `json:"penalty"`
}

func (f FixturesQuery) verb() string {
	return "GET"
}

func (f FixturesQuery) path() string {
	return "/fixtures"
}

// List ...
func (f FixturesQuery) List() ([]Fixtures, FixturesEnvelope, error) {
	result := FixturesEnvelope{}
	request(f.api, f, &result)

	return result.Response, result, nil
}
