package nelly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type DeviceCommunicator struct {
	Handle        *pcap.Handle
	PacketSniffer *PacketStreamer
	PacketWriter  *PacketStreamer
}

// TODO: research different snap lengths
// The same default as tcpdump.
const defaultSnapLen = 262144

// TODO: split creation of packet reader and writer to different files
func BuildDeviceCommunicator(device string, writeErrorHandler PacketProcessor) (*DeviceCommunicator, error) {
	handle, err := pcap.OpenLive(device, defaultSnapLen, true, -1)
	if err != nil {
		return nil, err
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	sniffer := PacketStreamer{
		nil,
		ConvertChanToPointerChan[gopacket.Packet](packetSource.Packets()),
	}
	writer := PacketStreamer{
		buildDeviceWriterFilter(handle, writeErrorHandler),
		nil,
	}
	devCommunicator := DeviceCommunicator{
		handle,
		&sniffer,
		&writer,
	}
	return &devCommunicator, nil
}

func buildDeviceWriterFilter(handle *pcap.Handle, errorProcessor PacketProcessor) func(packet *gopacket.Packet) {
	return func(packet *gopacket.Packet) {
		rawBytes, err := serializePacket(packet)
		if err != nil {
			sugar.Warn("Failed to serialize packet")
			if errorProcessor != nil {
				errorProcessor(packet)
			}
		}
		err = handle.WritePacketData(rawBytes)
		if err != nil {
			sugar.Warn("Failed to send packet to device")
			if errorProcessor != nil {
				errorProcessor(packet)
			}
		}
	}
}

func serializePacket(packet *gopacket.Packet) ([]byte, error) {
	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializePacket(buffer, gopacket.SerializeOptions{
		FixLengths:       false,
		ComputeChecksums: false,
	}, *packet)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
