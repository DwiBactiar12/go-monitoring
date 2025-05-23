package usecase

import (
	"context"
	"fmt"
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"time"

	"github.com/google/uuid"
)

type MonitoringUsecase struct {
	monitoringRepo iface.MonitoringRepository
}

func NewMonitoringUsecase(repo iface.MonitoringRepository) iface.MonitoringUseCase {
	return &MonitoringUsecase{
		monitoringRepo: repo,
	}
}   

// Store data monitoring ke InfluxDB
func (uc *MonitoringUsecase) StoreMonitoringData(ctx context.Context, data *entity.MonitoringData) error {
	if data.DeviceID == uuid.Nil {
		return fmt.Errorf("device_id tidak boleh kosong")
	}

	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}

	return uc.monitoringRepo.Store(ctx, data)
}

// Ambil data monitoring berdasarkan device_id dengan filter waktu dan limit data
func (uc *MonitoringUsecase) GetMonitoringDataByDevice(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time, limit int) ([]*entity.MonitoringData, error) {
	if deviceID == uuid.Nil {
		return nil, fmt.Errorf("device_id tidak boleh kosong")
	}
	if limit <= 0 {
		limit = 100 // default limit
	}
	if endTime.Before(startTime) {
		return nil, fmt.Errorf("endTime harus lebih besar dari startTime")
	}

	return uc.monitoringRepo.GetByDeviceID(ctx, deviceID, startTime, endTime, limit)
}

// Ambil data monitoring terbaru berdasarkan device_id
func (uc *MonitoringUsecase) GetLatestMonitoringData(ctx context.Context, deviceID uuid.UUID) (*entity.MonitoringData, error) {
	if deviceID == uuid.Nil {
		return nil, fmt.Errorf("device_id tidak boleh kosong")
	}

	return uc.monitoringRepo.GetLatestByDeviceID(ctx, deviceID)
}

// Ambil statistik monitoring perangkat
func (uc *MonitoringUsecase) GetMonitoringStats(ctx context.Context, deviceID uuid.UUID, startTime, endTime time.Time) (*entity.MonitoringStats, error) {
	if deviceID == uuid.Nil {
		return nil, fmt.Errorf("device_id tidak boleh kosong")
	}
	if endTime.Before(startTime) {
		return nil, fmt.Errorf("endTime harus lebih besar dari startTime")
	}

	return uc.monitoringRepo.GetStats(ctx, deviceID, startTime, endTime)
}

// Hapus data lama berdasarkan periode retensi
func (uc *MonitoringUsecase) DeleteOldMonitoringData(ctx context.Context, retentionPeriod time.Duration) error {
	if retentionPeriod <= 0 {
		return fmt.Errorf("retentionPeriod harus lebih dari 0")
	}

	return uc.monitoringRepo.DeleteOldData(ctx, retentionPeriod)
}
