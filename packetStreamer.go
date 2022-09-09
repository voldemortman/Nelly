package nelly

import (
	"github.com/google/gopacket"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

type PacketProcessor func(*gopacket.Packet)

type PacketStreamer struct {
	processor PacketProcessor
	source    chan *gopacket.Packet
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

func (stream *PacketStreamer) Start(quit chan bool) chan *gopacket.Packet {
	output := make(chan *gopacket.Packet)
	go func() {
		hasQuitBeenCalled := false
		for !hasQuitBeenCalled {
			select {
			case packet := <-stream.source:
				if stream.processor != nil {
					stream.processor(packet)
				}
				output <- packet
			case <-quit:
				hasQuitBeenCalled = true
				sugar.Debug("Quit has been called")
			}
		}
	}()
	return output
}

func InitializePacketStreamer(source chan *gopacket.Packet) *PacketStreamer {
	sugar = InitializeLogger()
	return &PacketStreamer{
		source:    source,
		processor: nil,
	}
}
