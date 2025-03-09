package device

import (
	"time"

	"golang-repo-pattern/internal/pkg/entity/device"
)

type (
	Servicer interface {
		CreateDevice(deviceData *device.Device) (int32, error)
		GetDeviceByID(id int32) (*device.Device, error)
		GetDevices(brand string, state *device.State) ([]device.Device, error)
		UpdateDevice(deviceData *device.Device) error
		DeleteDeviceByID(id int32) error
	}

	ServiceParams struct {
		Repo Repository
	}

	Service struct {
		repo Repository
	}
)

func NewService(params ServiceParams) Servicer {
	return &Service{
		repo: params.Repo,
	}
}

func (s *Service) CreateDevice(deviceData *device.Device) (int32, error) {
	deviceData.CreationTime = time.Now()

	result := s.repo.CreateDevice(deviceData)

	return deviceData.ID, result
}

func (s *Service) GetDeviceByID(id int32) (*device.Device, error) {
	device := &device.Device{}

	if err := s.repo.GetDeviceByID(device, id); err != nil {
		return nil, err
	}

	return device, nil
}

func (s *Service) GetDevices(brand string, state *device.State) ([]device.Device, error) {
	devices := []device.Device{}

	if err := s.repo.GetDevices(&devices, brand, state); err != nil {
		return nil, err
	}

	return devices, nil
}

func (s *Service) UpdateDevice(deviceData *device.Device) error {
	deviceData.CreationTime = time.Time{}

	d := &device.Device{}
	if err := s.repo.GetDeviceByID(d, deviceData.ID); err != nil {
		return err
	}

	if d.State.String() == device.InUse.String() &&
		(deviceData.Name != "" || deviceData.Brand != "") {
		return ErrUpdateDeviceInUse
	}

	return s.repo.UpdateDevice(deviceData)
}

func (s *Service) DeleteDeviceByID(id int32) error {
	d := &device.Device{}
	if err := s.repo.GetDeviceByID(d, id); err != nil {
		return err
	}

	if d.State.String() == device.InUse.String() {
		return ErrDeleteDeviceInUse
	}

	return s.repo.DeleteDeviceByID(id)
}
