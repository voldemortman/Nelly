package nelly

import "github.com/google/gopacket"

type HubDeviceRetriever struct {
	// ContinueWith: google how to create compareable struct, and then create type Device with string and expiration.
	//							 Change devices to be array of Device, and make get devices not return the source mac.
	devices []string
}

func (hub *HubDeviceRetriever) AddDevice(device string) {
	// TODO: find out how to get index of array
	if hub.devices.indexOf(device) != -1 {
		hub.devices = append(hub.devices, device)
	} else {
		sugar.Warnf("Tried to add device that already exists to hub. Device %s", device)
	}
}

func (hub *HubDeviceRetriever) RemoveDevice(device string) {
	// TODO: find out how to get index of array
	if hub.devices.indexOf(device) != -1 {
		// TODO: find out how to remove value from array
		hub.devices.remove(device)
	} else {
		sugar.Warnf("Tried to remove device that doesn't exits on hub. Device %s", device)
	}
}

func (hub *HubDeviceRetriever) GetAppropriateDevices(packet *gopacket.Packet) []string {
	return hub.devices
}
