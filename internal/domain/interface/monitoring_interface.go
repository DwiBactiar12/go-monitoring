package iface

import (
	"context"
	"monitoring/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type MonitoringRepository interface {
	Store(ctx context.Context, data *entity.MonitoringData) error
	GetByDeviceID(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time, limit int) ([]*entity.MonitoringData, error)
	GetLatestByDeviceID(ctx context.Context, deviceID uuid.UUID) (*entity.MonitoringData, error)
	GetStats(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time) (*entity.MonitoringStats, error)
	DeleteOldData(ctx context.Context, retentionPeriod time.Duration) error
}

type MonitoringUseCase interface {
	StoreMonitoringData(ctx context.Context, data *entity.MonitoringData) error
	GetMonitoringDataByDevice(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time, limit int) ([]*entity.MonitoringData, error)
	GetLatestMonitoringData(ctx context.Context, deviceID uuid.UUID) (*entity.MonitoringData, error)
	GetMonitoringStats(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time) (*entity.MonitoringStats, error)
	DeleteOldMonitoringData(ctx context.Context, retentionPeriod time.Duration) error
}
