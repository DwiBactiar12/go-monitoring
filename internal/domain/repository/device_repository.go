package repository

import (
	"context"
	"fmt"
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) iface.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *entity.Device) error {
	if err := r.db.WithContext(ctx).Create(device).Error; err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	return nil
}

func (r *deviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
	var device entity.Device
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&device).Error; err != nil {
		return nil, fmt.Errorf("failed to get device by id: %w", err)
	}
	return &device, nil
}

func (r *deviceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Device, error) {
	var devices []*entity.Device
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to get devices by user id: %w", err)
	}
	return devices, nil
}

func (r *deviceRepository) GetByMacAddress(ctx context.Context, macAddress string) (*entity.Device, error) {
	var device entity.Device
	if err := r.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(&device).Error; err != nil {
		return nil, fmt.Errorf("failed to get device by mac address: %w", err)
	}
	return &device, nil
}

func (r *deviceRepository) Update(ctx context.Context, device *entity.Device) error {
	if err := r.db.WithContext(ctx).Save(device).Error; err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}
	return nil
}

func (r *deviceRepository) UpdateOnlineStatus(ctx context.Context, deviceID uuid.UUID, isOnline bool) error {
	updates := map[string]interface{}{
		"is_online": isOnline,
		"last_seen": time.Now(),
	}

	if err := r.db.WithContext(ctx).Model(&entity.Device{}).Where("id = ?", deviceID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update device online status: %w", err)
	}
	return nil
}

func (r *deviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Device{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}

func (r *deviceRepository) List(ctx context.Context, limit, offset int) ([]*entity.Device, error) {
	var devices []*entity.Device
	if err := r.db.WithContext(ctx).Preload("User").Limit(limit).Offset(offset).Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	return devices, nil
}
