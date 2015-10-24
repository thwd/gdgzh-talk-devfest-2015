package main

import (
	"fmt"
	"net/http"
	"time"

	"go_meetup_zurich_fall_2015/demo/cityapi"
	"go_meetup_zurich_fall_2015/demo/weatherapi"

	"golang.org/x/net/context"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cityQuery, err := city.FromRequest(r) // HL
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200) // HL
	defer cancel()

	ctx = city.NewContext(ctx, cityQuery) // HL

	result, err := weather.Query(ctx) // HL
	if err != nil {
		http.Error(w, "Error retrieving weather", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Current temp is: %d\nForecast is: %d", result.Temperature, result.Forecast)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}
