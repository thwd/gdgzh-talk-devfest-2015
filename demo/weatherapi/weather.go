package weather

import (
	"encoding/json"
	"log"
	"net/http"

	"go_meetup_zurich_fall_2015/demo/cityapi"

	"golang.org/x/net/context"
)

const serviceBaseUrl = "http://localhost:20000/"

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }() // HL
	select {
	case <-ctx.Done(): // HL
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c: // HL
		return err
	}
}

type resultWrapper struct {
	Value int
	Err   error
}

func handleResponse(resp *http.Response, err error, res *int) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Temp int `json:"temp"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	*res = data.Temp
	return nil
}

func createRequest(reqType string, city string) *http.Request {
	req, err := http.NewRequest("GET", serviceBaseUrl+reqType, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Set("q", city)
	req.URL.RawQuery = q.Encode()

	return req
}

func getImpl(ctx context.Context, reqType string, out chan resultWrapper) {
	ci, _ := city.FromContext(ctx) // HL
	request := createRequest(reqType, ci)
	res := 0

	err := httpDo(ctx, request, func(resp *http.Response, err error) error { // HL
		return handleResponse(resp, err, &res)
	})

	out <- resultWrapper{res, err} // HL
}

func getForecast(ctx context.Context, out chan resultWrapper) {
	getImpl(ctx, "forecast", out)
}

func getTemp(ctx context.Context, out chan resultWrapper) {
	getImpl(ctx, "temp", out)
}

type QueryResult struct {
	Temperature int
	Forecast    int
}

func Query(ctx context.Context) (*QueryResult, error) {
	tempChan, forecastChan := make(chan resultWrapper, 1), make(chan resultWrapper, 1)
	go getTemp(ctx, tempChan)         // HL
	go getForecast(ctx, forecastChan) // HL
	temperature, forecast := resultWrapper{}, resultWrapper{}

	for i := 0; i < 2; i++ {
		select {
		case temperature = <-tempChan:
			if temperature.Err != nil {
				return nil, temperature.Err
			}
		case forecast = <-forecastChan:
			if forecast.Err != nil {
				return nil, forecast.Err
			}
		case <-ctx.Done(): // HL
			return nil, ctx.Err() // HL
		}
	}
	return &QueryResult{temperature.Value, forecast.Value}, nil
}
