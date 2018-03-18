package storage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/makizm/ssmartthings/lan-device-handler/routing/model"
	"github.com/manyminds/api2go"
)

// DeviceStorage stores all devices
type DeviceStorage struct {
	devices map[string]*model.Device
	idCount int
}

// NewDeviceStorage initializes the storage
func NewDeviceStorage() *DeviceStorage {
	return &DeviceStorage{make(map[string]*model.Device), 1}
}

// GetAll returns the device map (because we need the ID as key too)
func (s DeviceStorage) GetAll() map[string]*model.Device {
	return s.devices
}

// GetOne device
func (s DeviceStorage) GetOne(id string) (model.Device, error) {
	device, ok := s.devices[id]
	if ok {
		return *device, nil
	}
	errMessage := fmt.Sprintf("Device for id %s not found", id)
	return model.Device{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a device
func (s *DeviceStorage) Insert(c model.Device) string {
	id := fmt.Sprintf("%d", s.idCount)
	c.ID = id
	s.devices[id] = &c
	s.idCount++
	return id
}

// Delete one :(
func (s *DeviceStorage) Delete(id string) error {
	_, exists := s.devices[id]
	if !exists {
		return fmt.Errorf("Device with id %s does not exist", id)
	}
	delete(s.devices, id)

	return nil
}

// Update a device
func (s *DeviceStorage) Update(c model.Device) error {
	_, exists := s.devices[c.ID]
	if !exists {
		return fmt.Errorf("Device with id %s does not exist", c.ID)
	}
	s.devices[c.ID] = &c

	return nil
}
