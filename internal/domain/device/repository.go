package device

import (
	"gorm.io/gorm"

	"golang-repo-pattern/internal/pkg/entity/device"
)

type Repository interface {
	CreateDevice(deviceData *device.Device) error
	GetDeviceByID(target *device.Device, id int32) error
	GetDevices(target *[]device.Device, brand string, state *device.State) error
	UpdateDevice(deviceData *device.Device) error
	DeleteDeviceByID(id int32) error
	HandleError(res *gorm.DB) error
}
