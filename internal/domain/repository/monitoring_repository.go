package repository

import (
	"context"
	"fmt"
	"monitoring/config"
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"time"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type monitoringRepository struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
	bucket   string
	org      string
}

func NewMonitoringRepository(client influxdb2.Client, cfg *config.InfluxDBConfig) iface.MonitoringRepository {
	return &monitoringRepository{
		client:   client,
		writeAPI: client.WriteAPI(cfg.Org, cfg.Bucket),
		queryAPI: client.QueryAPI(cfg.Org),
		bucket:   cfg.Bucket,
		org:      cfg.Org,
	}
}

func (r *monitoringRepository) Store(ctx context.Context, data *entity.MonitoringData) error {
	point := write.NewPoint("device_monitoring",
		map[string]string{
			"device_id": data.DeviceID.String(),
		},
		map[string]interface{}{
			"cpu_usage":    data.CPUUsage,
			"memory_usage": data.MemoryUsage,
			"disk_usage":   data.DiskUsage,
			"temperature":  data.Temperature,
		},
		data.Timestamp)

	r.writeAPI.WritePoint(point)
	r.writeAPI.Flush()

	return nil
}

func (r *monitoringRepository) GetByDeviceID(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time, limit int) ([]*entity.MonitoringData, error) {
	query := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "device_monitoring")
            |> filter(fn: (r) => r.device_id == "%s")
            |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> limit(n: %d)
    `, r.bucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), deviceID.String(), limit)

	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query monitoring data: %w", err)
	}
	defer result.Close()

	var monitoringData []*entity.MonitoringData
	for result.Next() {
		record := result.Record()
		data := &entity.MonitoringData{
			DeviceID:  deviceID,
			Timestamp: record.Time(),
		}

		if cpu, ok := record.ValueByKey("cpu_usage").(float64); ok {
			data.CPUUsage = cpu
		}
		if memory, ok := record.ValueByKey("memory_usage").(float64); ok {
			data.MemoryUsage = memory
		}
		if disk, ok := record.ValueByKey("disk_usage").(float64); ok {
			data.DiskUsage = disk
		}
		if temp, ok := record.ValueByKey("temperature").(float64); ok {
			data.Temperature = temp
		}

		monitoringData = append(monitoringData, data)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("query error: %w", result.Err())
	}

	return monitoringData, nil
}

func (r *monitoringRepository) GetLatestByDeviceID(ctx context.Context, deviceID uuid.UUID) (*entity.MonitoringData, error) {
	query := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -1h)
            |> filter(fn: (r) => r._measurement == "device_monitoring")
            |> filter(fn: (r) => r.device_id == "%s")
            |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> last()
    `, r.bucket, deviceID.String())

	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest monitoring data: %w", err)
	}
	defer result.Close()

	if result.Next() {
		record := result.Record()
		data := &entity.MonitoringData{
			DeviceID:  deviceID,
			Timestamp: record.Time(),
		}

		if cpu, ok := record.ValueByKey("cpu_usage").(float64); ok {
			data.CPUUsage = cpu
		}
		if memory, ok := record.ValueByKey("memory_usage").(float64); ok {
			data.MemoryUsage = memory
		}
		if disk, ok := record.ValueByKey("disk_usage").(float64); ok {
			data.DiskUsage = disk
		}
		if temp, ok := record.ValueByKey("temperature").(float64); ok {
			data.Temperature = temp
		}

		return data, nil
	}

	return nil, fmt.Errorf("no data found")
}

func (r *monitoringRepository) GetStats(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time) (*entity.MonitoringStats, error) {
	query := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "device_monitoring")
            |> filter(fn: (r) => r.device_id == "%s")
            |> group(columns: ["_field"])
            |> aggregateWindow(every: 1h, fn: mean, createEmpty: false)
            |> yield(name: "mean")
    `, r.bucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), deviceID.String())

	// This is a simplified version - you might want to implement more detailed statistics
	result, err := r.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query monitoring stats: %w", err)
	}
	defer result.Close()

	stats := &entity.MonitoringStats{
		DeviceID: deviceID,
		Period:   fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
	}

	dataPoints := 0
	for result.Next() {
		dataPoints++
	}
	stats.DataPoints = dataPoints

	return stats, nil
}

func (r *monitoringRepository) DeleteOldData(ctx context.Context, retentionPeriod time.Duration) error {
	deleteTime := time.Now().Add(-retentionPeriod)

	deleteAPI := r.client.DeleteAPI()
	err := deleteAPI.DeleteWithName(ctx, r.org, r.bucket, deleteTime, time.Now(), "")
	if err != nil {
		return fmt.Errorf("failed to delete old data: %w", err)
	}

	return nil
}
