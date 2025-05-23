package iface

import (
	"context"
	"monitoring/internal/domain/entity"

	"github.com/google/uuid"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *entity.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Device, error)
	GetByMacAddress(ctx context.Context, macAddress string) (*entity.Device, error)
	Update(ctx context.Context, device *entity.Device) error
	UpdateOnlineStatus(ctx context.Context, deviceID uuid.UUID, isOnline bool) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Device, error)
}

type DeviceUseCase interface {
	Create(ctx context.Context, device *entity.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Device, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Device, error)
	GetByMacAddress(ctx context.Context, macAddress string) (*entity.Device, error)
	Update(ctx context.Context, device *entity.Device) error
	UpdateOnlineStatus(ctx context.Context, deviceID uuid.UUID, isOnline bool) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Device, error)
}
