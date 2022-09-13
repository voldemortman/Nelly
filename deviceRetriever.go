package nelly

import "github.com/google/gopacket"

type DeviceRetriever interface {
	AddDevice(string)
	RemoveDevice(string)
	GetAppropriateDevices(*gopacket.Packet) []string
}
