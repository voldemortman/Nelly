package nelly

import (
	"github.com/google/gopacket"
)

type PacketProcessor func(*gopacket.Packet)

type PacketStreamer struct {
	processor PacketProcessor
	Source    chan *gopacket.Packet
}

func (stream *PacketStreamer) AddProcessor(processorFunc PacketProcessor) *PacketStreamer {
	if stream.processor != nil {
		stream.processor = processorFunc
	} else {
		stream.processor = func(packet *gopacket.Packet) {
			stream.processor(packet)
			processorFunc(packet)
		}
	}
	return stream
}
