package apifootball

import (
	"testing"
)

var fixturesResponse = `
	{
		"get": "fixtures",
		"parameters": {
		"live": "all"
		},
		"errors": [],
		"results": 4,
		"response": [
		{
			"fixture": {
			"id": 239625,
			"venue": "Stade Municipal de Oued Zem, Oued Zem",
			"referee": null,
			"timezone": "UTC",
			"date": "2020-02-06T14:00:00+00:00",
			"timestamp": 1580997600,
			"periods": {
				"first": 1580997600,
				"second": null
			},
			"status": {
				"long": "Halftime",
				"short": "HT",
				"elapsed": 45
			}
			},
			"league": {
			"id": 200,
			"name": "Botola Pro",
			"country": "Morocco",
			"logo": "https://media.api-football.com/leagues/115.png",
			"flag": "https://media.api-football.com/flags/ma.svg",
			"season": 2019,
			"round": "Regular Season - 14"
			},
			"teams": {
			"home": {
				"id": 967,
				"name": "Rapide Oued ZEM",
				"logo": "https://media.api-football.com/teams/967.png"
			},
			"away": {
				"id": 968,
				"name": "Wydad AC",
				"logo": "https://media.api-football.com/teams/968.png"
			}
			},
			"goals": {
			"home": 0,
			"away": 1
			},
			"score": {
			"halftime": {
				"home": 0,
				"away": 1
			},
			"fulltime": {
				"home": null,
				"away": null
			},
			"extratime": {
				"home": null,
				"away": null
			},
			"penalty": {
				"home": null,
				"away": null
			}
			}
		}
		]
	}
`

func TestFixturesRequest(t *testing.T) {
	type FixturesTestEnvelope struct {
		Results  int        `json:"results"`
		Response []Fixtures `json:"response"`
	}
	server, api := startMockServer(t, []RequestHandler{RequestHandler{
		PathAndQuery:  "/fixtures",
		ContentToSend: fixturesResponse}})

	defer server.Close()

	query := FixturesQuery{}

	data := FixturesTestEnvelope{}
	response, _ := request(*api, query, &data)
	expect(t, response.StatusCode, 200)
	expect(t, data.Results, 4)

}

func TestList(t *testing.T) {

	params := map[string]interface{}{
		"id":   3711,
		"Live": "yes",
	}

	server, api := startMockServer(t, []RequestHandler{RequestHandler{
		PathAndQuery:  "/fixtures",
		ContentToSend: fixturesResponse}})

	defer server.Close()

	f, _, _ := api.Fixtures(params).List()

	expect(t, f[0].Fixture.ID, 239625)
	expect(t, f[0].League.Name, "Botola Pro")
	expect(t, f[0].Score.Halftime.Away, 1)
}
