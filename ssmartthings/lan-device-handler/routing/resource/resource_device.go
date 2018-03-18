package resource

import (
	"errors"
	"net/http"

	"github.com/makizm/ssmartthings/lan-device-handler/routing/model"
	"github.com/makizm/ssmartthings/lan-device-handler/routing/storage"

	"github.com/manyminds/api2go"
)

// DeviceResource for api2go routes
type DeviceResource struct {
	DeviceStorage *storage.DeviceStorage
}

// FindAll to satisfy api2go data source interface
func (d DeviceResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []model.Device
	devices := d.DeviceStorage.GetAll()

	for _, device := range devices {
		result = append(result, *device)
	}

	return &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the device with the given ID, otherwise an error
func (d DeviceResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	device, err := d.DeviceStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: device}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (d DeviceResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	device, ok := obj.(model.Device)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := d.DeviceStorage.Insert(device)
	device.ID = id

	return &Response{Res: device, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (d DeviceResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := d.DeviceStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

// Update stores all changes on the device
func (d DeviceResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	device, ok := obj.(model.Device)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := d.DeviceStorage.Update(device)
	return &Response{Res: device, Code: http.StatusNoContent}, err
}

// To do
// PaginatedFindAll can be used to load devices in chunks
// func (s DeviceResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) { }
