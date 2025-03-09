package device

import (
	"gorm.io/gorm"

	"golang-repo-pattern/internal/pkg/entity/device"
)

type Repository interface {
	CreateDevice(deviceData *device.Device) (int32, error)
	GetDeviceByID(target *device.Device, id int32) (*device.Device, error)
	GetDevices(target []device.Device, brand string, state *device.State) ([]device.Device, error)
	UpdateDevice(deviceData *device.Device) error
	DeleteDeviceByID(id int32) error
	HandleError(res *gorm.DB) error
}
