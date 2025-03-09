package device

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"golang-repo-pattern/internal/pkg/entity/device"
)

type repository struct {
	db           *gorm.DB
	defaultJoins []string
}

func NewRepository(db *gorm.DB, joins ...string) Repository {
	return repository{
		db:           db,
		defaultJoins: joins,
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

func (r repository) GetDevices(target []*device.Device, brand string, state *device.State) error {
	device := &device.Device{}

	if brand != "" {
		device.Brand = brand
	}

	if state != nil {
		device.State = *state
	}

	res := r.db.Find(target).Where(device)
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

func (r repository) DeleteDeviceByID(target *device.Device, id int32) error {
	res := r.db.Delete(target, id)
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
