# Go wrapper for Api Football

This is a Go api wrapper for the api football service. Currently Fixtures are the only thing that can be accessed.

https://www.api-football.com/

`
options := Options{
    key: "api-key"
}

client := NewClient(options)
params := map[string]interface{}{}

fixtures, r, err := api.Fixtures(params).List()
`