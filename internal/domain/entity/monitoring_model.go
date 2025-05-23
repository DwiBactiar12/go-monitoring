package entity

import (
	"time"

	"github.com/google/uuid"
)

type MonitoringData struct {
	DeviceID    uuid.UUID `json:"device_id"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

type MonitoringRequest struct {
	CPUUsage    float64 `json:"cpu_usage" validate:"required,gt=0"`
	MemoryUsage float64 `json:"memory_usage" validate:"required"`
	DiskUsage   float64 `json:"disk_usage" validate:"required"`
	Temperature float64 `json:"temperature" validate:"required"`
}

type MQTTMessage struct {
	DeviceID  string   `json:"device_id"`
	Type      string   `json:"type"`
	Data      MQTTData `json:"data"`
	Timestamp int64    `json:"timestamp"`
}

type MQTTData struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Temperature float64 `json:"temperature"`
	IPAddress   string  `json:"ip_address"`
}

type MonitoringQueryParams struct {
	DeviceID  uuid.UUID `query:"device_id"`
	StartTime time.Time `query:"start_time"`
	EndTime   time.Time `query:"end_time"`
	Limit     int       `query:"limit"`
}

type MonitoringStats struct {
	DeviceID   uuid.UUID `json:"device_id"`
	AvgCPU     float64   `json:"avg_cpu"`
	AvgMemory  float64   `json:"avg_memory"`
	AvgDisk    float64   `json:"avg_disk"`
	AvgTemp    float64   `json:"avg_temperature"`
	MaxCPU     float64   `json:"max_cpu"`
	MaxMemory  float64   `json:"max_memory"`
	MaxDisk    float64   `json:"max_disk"`
	MaxTemp    float64   `json:"max_temperature"`
	DataPoints int       `json:"data_points"`
	Period     string    `json:"period"`
}
