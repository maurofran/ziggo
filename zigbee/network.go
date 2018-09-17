package zigbee

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
)

// NetworkListener is the interface implemented by objects who needs to be
// notified by network changes.
type NetworkListener interface {
	DeviceAdded(Device)
	DeviceUpdated(Device)
	DeviceRemoved(Device)
}

const defaultStateFilePath = "simple-network.json"

// Network is the ZigBee network state implementation.
type Network struct {
	devices     map[string]Device
	devicesMx   sync.RWMutex
	groups      map[uint32]GroupAddress
	groupsMx    sync.RWMutex
	listeners   []NetworkListener
	listenersMx sync.RWMutex
	reset       bool
	filePath    string
}

// NewNetworkState will create a new NetworkState instance.
func NewNetworkState(reset bool) *Network {
	return &Network{
		devices:   make(map[string]Device),
		groups:    make(map[uint32]GroupAddress),
		listeners: nil,
		reset:     reset,
		filePath:  defaultStateFilePath,
	}
}

// Startup will start the network.
func (n *Network) Startup() error {
	filePath := n.filePath
	_, err := os.Stat(filePath)
	if !n.reset && err == nil {
		log.Println("Loading network state.")
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return errors.Wrapf(err, "Unable to read content of file %s", filePath)
		}
		if err := json.Unmarshal(bytes, n); err != nil {
			return errors.Wrapf(err, "Unable to unmarshal network state from file %s", filePath)
		}
		log.Println("Loading network state done.")
	}
	return nil
}

// Shutdown will stop the network.
func (n *Network) Shutdown() error {
	log.Println("Saving network state.")
	bytes, err := json.Marshal(n)
	if err != nil {
		return errors.Wrapf(err, "Unable to marshal network state to file %s", n.filePath)
	}
	if err := ioutil.WriteFile(n.filePath, bytes, 0644); err != nil {
		return errors.Wrapf(err, "Unabel to write content to file %s", n.filePath)
	}
	log.Println("Saving network state done.")
	return nil
}

// AddGroup will add the group address to this network.
func (n *Network) AddGroup(address GroupAddress) {
	n.groupsMx.Lock()
	defer n.groupsMx.Unlock()
	n.groups[address.GroupID] = address
}

// UpdateGroup will update the group address in this network.
func (n *Network) UpdateGroup(address GroupAddress) {
	n.groupsMx.Lock()
	defer n.groupsMx.Unlock()
	n.groups[address.GroupID] = address
}

// RemoveGroup will remove a group address from this network.
func (n *Network) RemoveGroup(address GroupAddress) {
	n.groupsMx.Lock()
	defer n.groupsMx.Unlock()
	delete(n.groups, address.GroupID)
}

// Group will retrieve the group address for supplied group id. The bool value is false if group address was not found.
func (n *Network) Group(groupID uint32) (GroupAddress, bool) {
	n.groupsMx.RLock()
	defer n.groupsMx.RUnlock()
	address, ok := n.groups[groupID]
	return address, ok
}

// Groups returns a copy of group addresses.
func (n *Network) Groups() []GroupAddress {
	n.groupsMx.RLock()
	defer n.groupsMx.RUnlock()
	var result []GroupAddress
	for _, address := range n.groups {
		result = append(result, address)
	}
	return result
}

// AddDevice will add a new device to network.
func (n *Network) AddDevice(device Device) {
	n.devicesMx.Lock()
	defer n.devicesMx.Unlock()
	n.devices[device.NetworkAddress.String()] = device
	n.listenersMx.RLock()
	defer n.listenersMx.RUnlock()
	for _, listener := range n.listeners {
		listener.DeviceAdded(device)
	}
}

// UpdateDevice will update an existing device.
func (n *Network) UpdateDevice(device Device) {
	n.devicesMx.Lock()
	defer n.devicesMx.Unlock()
	n.devices[device.NetworkAddress.String()] = device
	n.listenersMx.RLock()
	defer n.listenersMx.RUnlock()
	for _, listener := range n.listeners {
		listener.DeviceUpdated(device)
	}
}

// RemoveDevice will remove the device from network.
func (n *Network) RemoveDevice(device Device) {
	n.devicesMx.Lock()
	defer n.devicesMx.Unlock()
	delete(n.devices, device.NetworkAddress.String())
	n.listenersMx.RLock()
	defer n.listenersMx.RUnlock()
	for _, listener := range n.listeners {
		listener.DeviceRemoved(device)
	}
}

// Device will retrieve a device for supplied address. The bool value is false if no device is found.
func (n *Network) Device(address Address) (Device, bool) {
	if address.IsGroup() {
		return Device{}, false
	}
	n.devicesMx.RLock()
	defer n.devicesMx.RUnlock()
	device, ok := n.devices[address.String()]
	return device, ok
}

// Devices will retrieve a slices of all devices.
func (n *Network) Devices() []Device {
	n.devicesMx.RLock()
	defer n.devicesMx.RUnlock()
	var result []Device
	for _, device := range n.devices {
		result = append(result, device)
	}
	return result
}

// AddNetworkListener will add a network listener.
func (n *Network) AddNetworkListener(listener NetworkListener) {
	n.listenersMx.Lock()
	defer n.listenersMx.Unlock()
	for _, l := range n.listeners {
		if l == listener {
			return
		}
	}
	n.listeners = append(n.listeners, listener)
}

// RemoveNetworkListener will remove a network listener.
func (n *Network) RemoveNetworkListener(listener NetworkListener) {
	n.listenersMx.Lock()
	defer n.listenersMx.Unlock()
	for i, l := range n.listeners {
		if l == listener {
			n.listeners[i] = n.listeners[len(n.listeners)-1]
			n.listeners[len(n.listeners)-1] = nil
			n.listeners = n.listeners[:len(n.listeners)-1]
			return
		}
	}
}

type serializedNetwork struct {
	Devices []Device       `json:"devices"`
	Groups  []GroupAddress `json:"groups"`
}

// MarshalJSON will implement custom JSON serialization.
func (n *Network) MarshalJSON() ([]byte, error) {
	n.devicesMx.RLock()
	n.groupsMx.RLock()
	defer n.devicesMx.RUnlock()
	defer n.groupsMx.RUnlock()
	// Network state is a serialization of an array of devices and groups
	state := &serializedNetwork{
		Devices: n.Devices(),
		Groups:  n.Groups(),
	}
	return json.Marshal(state)
}

// UnmarshalJSON will implement custom JSON deserialization.
func (n *Network) UnmarshalJSON(data []byte) error {
	var state serializedNetwork
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	n.devicesMx.Lock()
	n.groupsMx.Lock()
	defer n.devicesMx.Unlock()
	defer n.groupsMx.Unlock()
	for _, device := range state.Devices {
		n.devices[device.NetworkAddress.String()] = device
	}
	for _, group := range state.Groups {
		n.groups[group.GroupID] = group
	}
	return nil
}
