package nelly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// TODO: research different snap lengths
// The same default as tcpdump.
const defaultSnapLen = 262144

func BuildDeviceHandleCommunication(device string, writeErrorHandler PacketFilter) (*DeviceCommunicator, error) {
	handle, err := pcap.OpenLive(device, defaultSnapLen, true, -1)
	if err != nil {
		return nil, err
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	sniffer := PacketStreamer{
		nil,
		ConvertChanToPointerChan[gopacket.Packet](packetSource.Packets()),
		nil,
	}
	writer := PacketStreamer{
		buildDeviceWriterFilter(handle, writeErrorHandler),
		nil,
		nil,
	}
	devCommunicator := DeviceCommunicator{
		handle,
		&sniffer,
		&writer,
	}
	return &devCommunicator, nil
}

func buildDeviceWriterFilter(handle *pcap.Handle, errorProcessor PacketFilter) *PacketFilter {
	var writerFilter PacketFilter
	writerFilter = func(packet *gopacket.Packet) {
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
	return &writerFilter
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
