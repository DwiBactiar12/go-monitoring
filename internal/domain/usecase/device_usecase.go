package usecase

import (
	"context"
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"monitoring/pkg/db"

	"github.com/google/uuid"
)

type deviceUsecase struct {
	deviceRepo iface.DeviceRepository
	cache0     *db.Client
}

func NewDeviceUsecase(deviceRepo iface.DeviceRepository, cache *db.Client) iface.DeviceUseCase {
	return &deviceUsecase{
		deviceRepo: deviceRepo,
		cache0:     cache,
	}
}

func (d *deviceUsecase) Create(ctx context.Context, device *entity.Device) error {
	return d.deviceRepo.Create(ctx, device)
}

func (d *deviceUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Device, error) {
    return d.deviceRepo.GetByID(ctx, id)
}

func (d *deviceUsecase) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Device, error) {
    return d.deviceRepo.GetByUserID(ctx, userID)
}

func (d *deviceUsecase) GetByMacAddress(ctx context.Context, macAddress string) (*entity.Device, error) {
    return d.deviceRepo.GetByMacAddress(ctx, macAddress)
}

func (d *deviceUsecase) Update(ctx context.Context, device *entity.Device) error {
    return d.deviceRepo.Update(ctx, device)
}

func (d *deviceUsecase) UpdateOnlineStatus(ctx context.Context, deviceID uuid.UUID, isOnline bool) error {
    return d.deviceRepo.UpdateOnlineStatus(ctx, deviceID, isOnline)
}

func (d *deviceUsecase) Delete(ctx context.Context, id uuid.UUID) error {
    return d.deviceRepo.Delete(ctx, id)
}

func (d *deviceUsecase) List(ctx context.Context, limit, offset int) ([]*entity.Device, error) {
    return d.deviceRepo.List(ctx, limit, offset)
}


