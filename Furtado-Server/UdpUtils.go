package furtado

import (
	"net"
)

func CreateUDPSocket(sourceAddress string, destAddress string) (*net.UDPConn, error) {
	udpAddrSource, err := net.ResolveUDPAddr("udp", sourceAddress)
	if err != nil {
		return nil, err
	}
	udpAddrDest, err := net.ResolveUDPAddr("udp", destAddress)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", udpAddrSource, udpAddrDest)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type UDPConnectionReadResult struct {
	buffer  *[]byte
	length  int
	UDPAddr *net.UDPAddr
	err     error
}

func ListenOnUDP(address string, readResults chan *UDPConnectionReadResult) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			buffer := make([]byte, 1024)
			length, address, err := conn.ReadFromUDP(buffer)
			result := &UDPConnectionReadResult{
				buffer:  &buffer,
				length:  length,
				UDPAddr: address,
				err:     err,
			}
			readResults <- result
		}
	}()

	return conn, nil
}
