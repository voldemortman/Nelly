package nelly

import (
	"github.com/google/gopacket"
)

type PacketFilter func(*gopacket.Packet)

type PacketFilterBuilder struct {
	filters []PacketFilter
}

func (filterBuilder *PacketFilterBuilder) AddFilter(filter PacketFilter) {
	if filterBuilder.filters == nil {
		filterBuilder.filters = []PacketFilter{filter}
	} else {
		filterBuilder.filters = append(filterBuilder.filters, filter)
	}
}

func (filterBuilder *PacketFilterBuilder) BuildFilter() *PacketFilter {
	var packetFilter PacketFilter
	for _, filter := range filterBuilder.filters {
		if packetFilter != nil {
			packetFilter = func(p *gopacket.Packet) {
				packetFilter(p)
				filter(p)
			}
		} else {
			packetFilter = filter
		}
	}
	return &packetFilter
}
