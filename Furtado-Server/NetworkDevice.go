package furtado

import (
	"github.com/google/gopacket/pcap"
)

type DeviceCommunicator struct {
	Handle *pcap.Handle
}
