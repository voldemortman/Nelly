package nelly

import (
	"errors"

	"github.com/google/gopacket"
)

type PacketStreamer struct {
	Filter *PacketFilter
	source chan *gopacket.Packet
	output chan *gopacket.Packet
}

func (stream *PacketStreamer) AddSource(source chan *gopacket.Packet) error {
	if stream.source == nil {
		stream.source = source
		return nil
	}
	return errors.New("packet streamer already has source")
}

func (stream *PacketStreamer) AddOutput(output chan *gopacket.Packet) error {
	if stream.output == nil {
		stream.output = output
		return nil
	}
	return errors.New("packet streamer already has output")
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
				if stream.output != nil {
					stream.output <- packet
				}
			case <-quit:
				hasQuitBeenCalled = true
				sugar.Debug("Quit has been called")
			}
		}
	}()
}
