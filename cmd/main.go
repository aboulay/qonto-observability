package main

import (
	"os"
	"fmt"
	"net/http"
	"time"
	
	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"qonto-observability/internal/domain"
	"qonto-observability/internal/outbound/openweathermap"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Batches []*domain.Batch
	ApiKey string
}

func setupConfiguration() (*Config, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	filePath := os.Getenv("CITY_FILE_PATH")
	cities := domain.RetrieveCitiesFromFile(filePath)

	batches := domain.GenerateBatches(cities)

	config := Config{
		Batches: batches,
		ApiKey: apiKey,
	}

	return &config, nil
}

func processBatch(batch *domain.Batch, apiKey string, metric *prometheus.GaugeVec) error {
	for _, city := range batch.Cities {
		weather, err := openweathermap.GetWeather(city, apiKey)
		if err != nil {
			return fmt.Errorf("an error occured when processing the city %s: %v", city.Name, err)
		}
		logrus.Debug(fmt.Sprintf("%s-%f", city.Name, weather.Main.Temp))
		metric.WithLabelValues(city.Name).Set(weather.Main.Temp)
	}
	return nil
}

func runBatchSystem(config *Config, metric *prometheus.GaugeVec) {
	for {
		for _, batch := range config.Batches {
			if err := processBatch(batch, config.ApiKey, metric); err != nil {
				logrus.Fatal(err)
			}
			logrus.Info("waiting 1min30 before the next batch")
			time.Sleep(1*time.Minute + 30*time.Second)
		}
		logrus.Info("every batch has been handled, come back at the beginning of the list")
	}
}

func setupMetric() *prometheus.GaugeVec {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "city_temp",
			Help: "temperature in Celsius",
		},
		[]string{"city"},
	)
	prometheus.MustRegister(metric)
	return metric
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	config, err := setupConfiguration()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("configuration initialized")

	metric := setupMetric()

	go runBatchSystem(config, metric)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}