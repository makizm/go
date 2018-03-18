package model

// Device is a generic child device
type Device struct {
	ID     string `json:"-"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	exists bool
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (d Device) GetID() string {
	return d.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (d *Device) SetID(id string) error {
	d.ID = id
	return nil
}
