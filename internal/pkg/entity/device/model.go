package device

import "time"

type Device struct {
	ID    int32  `gorm:"primarykey" json:"id"`
	Name  string `json:"name" validate:"required"`
	Brand string `json:"brand" validate:"required"`
	State State  `json:"state" validate:"required"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
