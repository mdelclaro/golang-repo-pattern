package device

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"golang-repo-pattern/internal/pkg/entity/device"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return repository{
		db: db,
	}
}

func (r repository) CreateDevice(deviceData *device.Device) error {
	res := r.db.Create(deviceData)
	return r.HandleError(res)
}

func (r repository) GetDeviceByID(target *device.Device, id int32) error {
	res := r.db.First(target, id)
	return r.HandleError(res)
}

func (r repository) GetDevices(target *[]device.Device, brand string, state *device.State) error {
	m := make(map[string]interface{})

	if brand != "" {
		m["brand"] = brand
	}

	if state != nil {
		m["state"] = *state
	}

	res := r.db.Where(m).Find(target)
	return r.HandleError(res)
}

func (r repository) UpdateDevice(deviceData *device.Device) error {
	res := r.db.
		Model(deviceData).
		Clauses(clause.Returning{}).
		Updates(deviceData)

	if res.RowsAffected == 0 {
		res.Error = fmt.Errorf("record not found")
	}

	return r.HandleError(res)
}

func (r repository) DeleteDeviceByID(id int32) error {
	res := r.db.Delete(&device.Device{}, id)
	if res.RowsAffected == 0 {
		res.Error = fmt.Errorf("record not found")
	}

	return r.HandleError(res)
}

func (r repository) HandleError(res *gorm.DB) error {
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		err := fmt.Errorf("%w", res.Error)
		return err
	}

	return nil
}
