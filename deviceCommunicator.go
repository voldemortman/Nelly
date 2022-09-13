package nelly

import (
	"errors"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type DeviceCommunicator struct {
	Handle        *pcap.Handle
	PacketSniffer *PacketStreamer
	PacketWriter  *PacketStreamer
}

type PacketFilter func(*gopacket.Packet)

type PacketStreamer struct {
	Filter *PacketFilter
	source chan *gopacket.Packet
}

func (stream *PacketStreamer) AddSource(source chan *gopacket.Packet) error {
	if stream.source == nil {
		stream.source = source
		return nil
	}
	return errors.New("packet streamer already has source")
}

func (stream *PacketStreamer) GetSource() chan *gopacket.Packet {
	return stream.source
}

func (stream *PacketStreamer) Start(quit chan bool) {
	go func() {
		hasQuitBeenCalled := false
		for !hasQuitBeenCalled {
			select {
			case packet := <-stream.source:
				if stream.Filter != nil && *stream.Filter != nil {
					(*stream.Filter)(packet)
				}
			case <-quit:
				hasQuitBeenCalled = true
				sugar.Debug("Quit has been called")
			}
		}
	}()
}
