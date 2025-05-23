package db

import (
	"context"
	"monitoring/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxDB(cfg *config.Config) (influxdb2.Client, error) {
	client := influxdb2.NewClient(cfg.InfluxDB.URL, cfg.InfluxDB.Token)

	// Test connection
	_, err := client.Health(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}
