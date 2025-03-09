package device_integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v3"
	"github.com/zeebo/assert"

	"golang-repo-pattern/internal/domain/device"
	"golang-repo-pattern/internal/infra/database"
	deviceentity "golang-repo-pattern/internal/pkg/entity/device"
)

var (
	app *fiber.App
	id  int32 = 1
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestDevice(t *testing.T) {
	app = fiber.New()

	db, gormDb, mock := database.StartDbMock(t)
	repo := device.NewRepository(gormDb)

	service := device.NewService(device.ServiceParams{Repo: repo})
	device.NewHttpHandler(app, service)

	defer db.Close()

	tests := []struct {
		name string

		route  string
		method string
		body   any

		expectedCode int
		expectedBody any

		mock func()
	}{
		{
			name:         "[Success] - Test Create Device",
			route:        device.BaseRoute,
			method:       "POST",
			body:         deviceentity.Device{Name: "Device1", Brand: "Brand1", State: deviceentity.Available},
			expectedCode: http.StatusOK,
			expectedBody: deviceentity.Device{ID: 1, Name: "Device1", Brand: "Brand1", State: deviceentity.Available},
			mock: func() {
				row := sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Device1", "Brand1", deviceentity.Available, time.Time{})

				expectedSQL := "INSERT INTO \"devices\" (.+) VALUES (.+)"

				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()
			},
		},
		{
			name:         "[Error] - Test Create Device with Invalid Data",
			route:        device.BaseRoute,
			method:       "POST",
			body:         deviceentity.Device{Name: "", Brand: "", State: deviceentity.Available},
			expectedCode: http.StatusBadRequest,
			expectedBody: fiber.Map{"error": "invalid field(s): Name, Brand"},
			mock: func() {
				// No database interaction needed for validation error
			},
		},
		{
			name:         "[Success] - Test Get Device By ID",
			route:        fmt.Sprintf(device.BaseRoute+"/%d", id),
			method:       "GET",
			expectedCode: http.StatusOK,
			expectedBody: deviceentity.Device{ID: 1, Name: "Device1", Brand: "Brand1", State: deviceentity.Available, CreationTime: time.Time{}},
			mock: func() {
				row := sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Device1", "Brand1", deviceentity.Available, time.Time{})

				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
			},
		},
		{
			name:         "[Error] - Test Get Device By Non-Existent ID",
			route:        fmt.Sprintf(device.BaseRoute+"/%d", id),
			method:       "GET",
			expectedCode: http.StatusNotFound,
			expectedBody: fiber.Map{"error": "device not found"},
			mock: func() {
				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnError(fmt.Errorf("device not found"))
			},
		},
		{
			name:         "[Success] - Test Get Devices",
			route:        device.BaseRoute + "?brand=Brand1",
			method:       "GET",
			expectedCode: http.StatusOK,
			expectedBody: []deviceentity.Device{{ID: 1, Name: "Device1", Brand: "Brand1", State: deviceentity.Available}},
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Device1", "Brand1", deviceentity.Available, time.Time{})

				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnRows(rows)
			},
		},
		{
			name:         "[Error] - Test Get Devices with Database Error",
			route:        device.BaseRoute + "?brand=Brand1",
			method:       "GET",
			expectedCode: http.StatusInternalServerError,
			expectedBody: fiber.Map{"error": "internal server error"},
			mock: func() {
				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnError(fmt.Errorf("internal server error"))
			},
		},
		{
			name:         "[Success] - Test Update Device",
			route:        fmt.Sprintf(device.BaseRoute+"/%d", id),
			method:       "PUT",
			body:         deviceentity.Device{Name: "Updated Device1", Brand: "Updated Brand1", State: deviceentity.Available},
			expectedCode: http.StatusOK,
			expectedBody: deviceentity.Device{ID: 1, Name: "Updated Device1", Brand: "Updated Brand1", State: deviceentity.Available},
			mock: func() {
				row := sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Device1", "Brand1", deviceentity.Available, time.Time{})

				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnRows(row).WithArgs(sqlmock.AnyArg(), id)

				row = sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Updated Device1", "Updated Brand1", deviceentity.Available, time.Time{})

				expectedSQL = "UPDATE \"devices\" SET (.+) WHERE (.+)"

				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
				mock.ExpectCommit()
			},
		},
		{
			name:         "[Error] - Test Update Device Name In Use",
			route:        fmt.Sprintf(device.BaseRoute+"/%d", id),
			method:       "PUT",
			body:         deviceentity.Device{Name: "Updated Device1", Brand: "Updated Brand1", State: deviceentity.Available},
			expectedCode: http.StatusBadRequest,
			expectedBody: fiber.Map{"error": "name/brand cannot be updated when device is in use"},
			mock: func() {
				row := sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Device1", "Brand1", deviceentity.InUse, time.Time{})

				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
			},
		},
		{
			name:         "[Success] - Test Delete Device",
			route:        fmt.Sprintf(device.BaseRoute+"/%d", id),
			method:       "DELETE",
			expectedCode: 204,
			expectedBody: "",
			mock: func() {
				device := sqlmock.NewRows([]string{
					"id", "name", "brand", "status", "creation_time",
				}).
					AddRow(id, "name", "brand", deviceentity.InUse, time.Time{})

				expectedSQL := "SELECT (.+) FROM \"devices\""
				mock.ExpectQuery(expectedSQL).WillReturnRows(device)

				expectedSQL = "DELETE FROM \"devices\" .+"

				mock.ExpectBegin()
				mock.ExpectExec(expectedSQL).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:         "[Error] - Test Delete Device In Use",
			route:        fmt.Sprintf(device.BaseRoute+"/%d", id),
			method:       "DELETE",
			expectedCode: http.StatusBadRequest,
			expectedBody: fiber.Map{"error": "cannot delete device in use"},
			mock: func() {
				row := sqlmock.NewRows([]string{
					"id", "name", "brand", "state", "creation_time",
				}).
					AddRow(id, "Device1", "Brand1", deviceentity.InUse, time.Time{})

				expectedSQL := "SELECT (.+) FROM \"devices\" WHERE (.+)"

				mock.ExpectQuery(expectedSQL).WillReturnRows(row)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			reqBody, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			bodyReader := bytes.NewReader(reqBody)

			req, _ := http.NewRequest(
				tt.method,
				tt.route,
				bodyReader,
			)

			res, err := app.Test(req, fiber.TestConfig{Timeout: -1})
			assert.NoError(t, err)

			body, _ := io.ReadAll(res.Body)
			parsedBody, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			if string(parsedBody) != "\"\"" {
				assert.Equal(t, string(parsedBody), string(body))
			}

			assert.Equal(t, tt.expectedCode, res.StatusCode)
		})
	}
}
