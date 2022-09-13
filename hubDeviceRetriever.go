package nelly

import "github.com/google/gopacket"

type HubDeviceRetriever struct {
	// ContinueWith: google how to create compareable struct, and then create type Device with string and expiration.
	//							 Change devices to be array of Device, and make get devices not return the source mac.
	devices []string
}

func (hub *HubDeviceRetriever) AddDevice(device string) {
	if hub.getDeviceIndex(device) == -1 {
		hub.devices = append(hub.devices, device)
	} else {
		sugar.Warnf("Tried to add device that already exists to hub. Device %s", device)
	}
}

func (hub *HubDeviceRetriever) RemoveDevice(device string) {
	if deviceIndex := hub.getDeviceIndex(device); deviceIndex != -1 {
		hub.devices[deviceIndex] = hub.devices[len(hub.devices)-1]
		hub.devices[len(hub.devices)-1] = ""
		hub.devices = hub.devices[:len(hub.devices)-1]
	} else {
		sugar.Warnf("Tried to remove device that doesn't exits on hub. Device %s", device)
	}
}

func (hub *HubDeviceRetriever) GetAppropriateDevices(packet *gopacket.Packet) []string {
	return hub.devices
}

func (hub *HubDeviceRetriever) getDeviceIndex(deviceToCheck string) int {
	for i, device := range hub.devices {
		if deviceToCheck == device {
			return i
		}
	}
	return -1
}
