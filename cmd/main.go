package main

import (
	"os"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"qonto-observability/internal/domain"
	"qonto-observability/internal/outbound/openweathermap"
)

func setup() (*Config, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	cities := domain.RetrieveCitiesFromFile("./cities.txt")

	batches := domain.GenerateBatches(cities)
	// for i, batch := range batches {
	// 	fmt.Printf("Batch %d:\n", i+1)
	// 	for _, city := range batch.Cities {
	// 		fmt.Printf("%s,%s", city.Name, city.Country)
	// 	}
	// }

	// response, err := openweathermap.GetWeather(cities[0], apiKey)
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// fmt.Printf("%f", response.Main.Temp)

	config := Config{
		Batches: batches
		ApiKey: apiKey
	}

	return &config, nil
}

func main() {
	config, error := setupConfiguration()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}