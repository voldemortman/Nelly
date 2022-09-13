package nelly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// TODO: research different snap lengths
// The same default as tcpdump.
const defaultSnapLen = 262144

func BuildDeviceHandleCommunication(device string, writeErrorHandler PacketFilter, macAddressTable AddressExpiration) (*DeviceCommunicator, error) {
	handle, err := pcap.OpenLive(device, defaultSnapLen, true, -1)
	if err != nil {
		return nil, err
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	sniffer := PacketStreamer{
		buildMacLearnerFilter(nil),
		ConvertChanToPointerChan[gopacket.Packet](packetSource.Packets()),
	}
	// TODO: should writer actually be a streamer? Maybe it needs to be a seperate struct that instead of using chan it writes one packet at a time?
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

func buildMacLearnerFilter(macAddresses *AddressExpiration) *PacketFilter {
	var filter PacketFilter
	filter = func(packet *gopacket.Packet) {
		sourceMac := (*packet).LinkLayer().LinkFlow().Src()
		err := (*macAddresses).UpdateAddressTimeStamp(sourceMac)
		if err != nil {
			sugar.Error(err)
		}
	}
	return &filter
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
