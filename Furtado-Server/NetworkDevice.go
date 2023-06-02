package furtado

import (
	"errors"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type NetworkDevice struct {
	handle      *pcap.Handle
	IsRunning   bool
	stopChannel chan bool
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

func (device *NetworkDevice) SendToDevice(packetData []byte) {
	err := device.handle.WritePacketData(packetData)
	if err != nil {
		fmt.Println("Error while writing packet to device: ", err)
	}
}

func (device *NetworkDevice) CloseCommunication() {
	device.stopChannel <- true
	device.handle.Close()
}
