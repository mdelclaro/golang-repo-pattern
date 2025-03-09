package device

import "errors"

var (
	ErrUpdateDeviceInUse = errors.New("name/brand cannot be updated when device is in use")
	ErrDeleteDeviceInUse = errors.New("cannot delete device in use")
)
