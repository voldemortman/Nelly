package netUtils

import (
	"bytes"
	"encoding/gob"

	"github.com/google/gopacket"
)

func SerializePacket(packet *gopacket.Packet) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(packet)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeserializePacket(data []byte) (*gopacket.Packet, error) {
	var packet gopacket.Packet
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&packet)
	if err != nil {
		return nil, err
	}
	return &packet, nil
}
