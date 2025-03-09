package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	v "github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"

	"golang-repo-pattern/internal/pkg/common"
	"golang-repo-pattern/internal/pkg/entity/device"
)

const baseRoute = "/devices"

type httpHandler struct {
	service Servicer
}

func NewHttpHandler(app *fiber.App, service Servicer) {
	handler := &httpHandler{
		service: service,
	}

	app.Post(baseRoute, handler.createDevice)
	app.Get(baseRoute+"/:id", handler.getDeviceByID)
	app.Get(baseRoute, handler.getDevices)
	app.Put(baseRoute+"/:id", handler.updateDevice)
	app.Delete(baseRoute+"/:id", handler.deleteDeviceByID)
}

func (h *httpHandler) createDevice(c fiber.Ctx) error {
	device := &device.Device{}

	if err := json.Unmarshal(c.Body(), &device); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	validator := v.New()
	if err := validator.Struct(device); err != nil {
		invalid := ""

		var valErrs v.ValidationErrors
		if errors.As(err, &valErrs) {

			for i, err := range valErrs {
				field := err.Field()
				if i != 0 {
					field = ", " + field
				}

				invalid = invalid + field
			}
		}

		return c.Status(http.StatusBadRequest).JSON(common.BuildError(fmt.Errorf("invalid field(s): %s", invalid)))
	}

	id, err := h.service.CreateDevice(device)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	device.ID = id

	return c.JSON(device)
}

func (h *httpHandler) getDeviceByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.BuildError(err))
	}

	device, err := h.service.GetDeviceByID(int32(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	return c.JSON(device)
}

func (h *httpHandler) getDevices(c fiber.Ctx) error {
	brand := c.Query("brand")
	state := c.Query("state")

	var deviceState *device.State
	if state != "" {
		deviceState = device.StringToState(state)
	}

	devices, err := h.service.GetDevices(brand, deviceState)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	return c.JSON(devices)
}

func (h *httpHandler) updateDevice(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.BuildError(err))
	}

	device := &device.Device{}
	if err := json.Unmarshal(c.Body(), &device); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	validator := v.New()
	if err := validator.Struct(device); err != nil {
		invalid := ""

		var valErrs v.ValidationErrors
		if errors.As(err, &valErrs) {

			for i, err := range valErrs {
				field := err.Field()
				if i != 0 {
					field = ", " + field
				}

				invalid = invalid + field
			}
		}

		return c.Status(http.StatusBadRequest).JSON(common.BuildError(fmt.Errorf("invalid field(s): %s", invalid)))
	}

	device.ID = int32(id)
	if err := h.service.UpdateDevice(device); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	return c.JSON(device)
}

func (h *httpHandler) deleteDeviceByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.BuildError(err))
	}

	if err := h.service.DeleteDeviceByID(int32(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.BuildError(err))
	}

	return c.SendStatus(http.StatusNoContent)
}
