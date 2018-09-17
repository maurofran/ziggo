package zigbee

import "fmt"

// Address is a generic interface for a zigbee address.
type Address interface {
	fmt.Stringer
	IsGroup() bool
}

// DeviceAddress defines a unicast ZigBee address.
type DeviceAddress struct {
	NetworkAddress uint32 `json:"networkAddress"`
	Endpoint       uint32 `json:"endpoint"`
}

// IsGroup will identify this address as not be a group address.
func (a DeviceAddress) IsGroup() bool {
	return false
}

func (a DeviceAddress) String() string {
	return fmt.Sprintf("%d/%d", a.NetworkAddress, a.Endpoint)
}

// GroupAddress defines a group ZigBee address.
type GroupAddress struct {
	GroupID uint32 `json:"groupId"`
	Label   string `json:"label"`
}

// IsGroup will identify this address as to be a group address.
func (a GroupAddress) IsGroup() bool {
	return true
}

func (a GroupAddress) String() string {
	return fmt.Sprintf("%d/%s", a.GroupID, a.Label)
}
