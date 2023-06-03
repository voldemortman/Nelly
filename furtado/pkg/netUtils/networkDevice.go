package netUtils

import (
	"errors"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type NetworkDevice struct {
	handle      *pcap.Handle
	IsRunning   bool
	stopChannel chan struct{}
}

const (
	// The same default as tcpdump
	defaultSnapLen = 262144
	// read here for why -1: https://pkg.go.dev/github.com/google/gopacket@v1.1.19/pcap#hdr-PCAP_Timeouts
	timeout = -1
)

func CreateNetworkDevice(deviceName string) (*NetworkDevice, error) {
	handle, err := pcap.OpenLive(deviceName, defaultSnapLen, true, timeout)
	if err != nil {
		return nil, err
	}

	return &NetworkDevice{
		handle:      handle,
		IsRunning:   false,
		stopChannel: make(chan struct{}),
	}, nil
}

func (device *NetworkDevice) ReadFromDevice() (chan *gopacket.Packet, error) {
	if device.IsRunning {
		return nil, errors.New("device is already in use")
	}
	device.IsRunning = true

	packetChan := make(chan *gopacket.Packet)
	go func() {
		for {
			select {
			case <-device.stopChannel:
				close(packetChan)
				return
			default:
				packetData, _, err := device.handle.ReadPacketData()
				if err != nil {
					fmt.Println("Error while reading packet: ", err)
					return
				}
				packet := gopacket.NewPacket(packetData, device.handle.LinkType(), gopacket.Default)
				packetChan <- &packet
			}
		}
	}()
	return packetChan, nil
}

func (device *NetworkDevice) SendToDevice(packetData *[]byte) {
	err := device.handle.WritePacketData(*packetData)
	if err != nil {
		fmt.Println("Error while writing packet to device: ", err)
	}
}

func (device *NetworkDevice) CloseCommunication() {
	close(device.stopChannel)
	device.handle.Close()
}
