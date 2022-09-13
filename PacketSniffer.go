package nelly

import (
	"github.com/google/gopacket"
)

type PacketSniffer struct {
	filterBuilder *PacketFilterBuilder
	source        chan *gopacket.Packet
}

func (sniffer *PacketSniffer) StartSniffing(quit chan bool) chan *gopacket.Packet {
	filter := sniffer.filterBuilder.BuildFilter()
	output := make(chan *gopacket.Packet)
	go func() {
		hasQuitBeenCalled := false
		for !hasQuitBeenCalled {
			select {
			case packet := <-sniffer.source:
				if filter != nil {
					(*filter)(packet)
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
