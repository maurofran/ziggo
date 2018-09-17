package zigbee

import "fmt"

// Device will represent a zigbee device.
type Device struct {
	IEEEAddress      uint64        `json:"ieeeAddress"`
	NetworkAddress   DeviceAddress `json:"networkAddress"`
	ProfileID        uint32        `json:"profileId"`
	DeviceType       uint32        `json:"deviceType"`
	DeviceID         uint32        `json:"deviceId"`
	ManufacturerCode uint32        `json:"manufacturerCode"`
	DeviceVersion    uint32        `json:"deviceVersion"`
	InputClusterIds  []uint32      `json:"inputClusterIds"`
	OutputClusterIds []uint32      `json:"outputClusterIds"`
	Label            string        `json:"label"`
}

func (d Device) String() string {
	return fmt.Sprintf(
		"ZigBeeDevice label=%s, networkAddress=%s, ieeeAddress=%x, profileId=%d, deviceType=%d, deviceId=%d, manufacturerCode=%d, deviceVersion=%d, inputClusterIds=%v, outputClusterIds=%v",
		d.Label,
		d.NetworkAddress,
		d.IEEEAddress,
		d.ProfileID,
		d.DeviceType,
		d.DeviceID,
		d.ManufacturerCode,
		d.DeviceVersion,
		d.InputClusterIds,
		d.OutputClusterIds,
	)
}
