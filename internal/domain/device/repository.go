package device

import (
	"gorm.io/gorm"

	"golang-repo-pattern/internal/pkg/entity/device"
)

type Repository interface {
	CreateDevice(deviceData *device.Device) (int32, error)
	GetDeviceByID(id int32) (*device.Device, error)
	GetDevices(brand string, state *device.State) ([]device.Device, error)
	UpdateDevice(deviceData *device.Device) error
	DeleteDeviceByID(id int32) error
	HandleError(res *gorm.DB) error
}
