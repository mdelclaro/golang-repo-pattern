package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	devicemock "golang-repo-pattern/internal/domain/device/mock"
	"golang-repo-pattern/internal/pkg/entity/device"
)

func TestDevice(t *testing.T) {
	mockRepo := devicemock.NewMockRepository(gomock.NewController(t))
	deviceData := &device.Device{Name: "Device1", Brand: "Brand1"}

	t.Run("TestCreateDevice", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		mockRepo.EXPECT().CreateDevice(deviceData).Return(int32(1), nil)

		id, err := service.CreateDevice(deviceData)
		assert.NoError(t, err)
		assert.Equal(t, id, int32(1))
	})

	t.Run("TestGetDeviceByID", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		mockRepo.EXPECT().GetDeviceByID(&device.Device{}, int32(1)).Return(deviceData, nil)

		result, err := service.GetDeviceByID(1)
		assert.NoError(t, err)
		assert.Equal(t, result.Brand, deviceData.Brand)
		assert.Equal(t, result.Name, deviceData.Name)
	})

	t.Run("TestGetDevices", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		devicesData := []device.Device{{Name: "Device1", Brand: "Brand1"}}

		mockRepo.EXPECT().GetDevices([]device.Device{}, "Brand1", nil).Return(devicesData, nil)

		result, err := service.GetDevices("Brand1", nil)
		assert.NoError(t, err)
		assert.Equal(t, devicesData, result)
	})

	t.Run("TestUpdateDevice", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		deviceData := &device.Device{ID: 1, Name: "Device1"}
		existingDevice := &device.Device{ID: 1, State: device.Available}
		mockRepo.EXPECT().GetDeviceByID(&device.Device{}, deviceData.ID).Return(existingDevice, nil)
		mockRepo.EXPECT().UpdateDevice(deviceData).Return(nil)

		err := service.UpdateDevice(deviceData)
		assert.NoError(t, err)
	})

	t.Run("TestUpdateDeviceInUse", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		deviceData := &device.Device{ID: 1, Name: "Device1"}
		existingDevice := &device.Device{ID: 1, State: device.InUse}
		mockRepo.EXPECT().GetDeviceByID(&device.Device{}, deviceData.ID).Return(existingDevice, nil)

		err := service.UpdateDevice(deviceData)
		assert.Error(t, err)
		assert.Equal(t, ErrUpdateDeviceInUse, err)
	})

	t.Run("TestDeleteDeviceByID", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		deviceData := &device.Device{ID: 1, State: device.Available}
		mockRepo.EXPECT().GetDeviceByID(&device.Device{}, int32(1)).Return(deviceData, nil)
		mockRepo.EXPECT().DeleteDeviceByID(int32(1)).Return(nil)

		err := service.DeleteDeviceByID(1)
		assert.NoError(t, err)
	})

	t.Run("TestDeleteDeviceInUse", func(t *testing.T) {
		t.Parallel()
		service := NewService(ServiceParams{Repo: mockRepo})

		deviceData := &device.Device{ID: 1, State: device.InUse}
		mockRepo.EXPECT().GetDeviceByID(&device.Device{}, int32(1)).Return(deviceData, nil)

		err := service.DeleteDeviceByID(1)
		assert.Error(t, err)
		assert.Equal(t, ErrDeleteDeviceInUse, err)
	})
}
