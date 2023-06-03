package netUtils

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
	Buffer  *[]byte
	Length  int
	UDPAddr *net.UDPAddr
	Err     error
}

// We return conn so that the user can close it. We do this on top of a quit channel is because we dont
// want a timeout on the read. So we want the user to close the connection. However, we don't want to start
// checking the contents of the error so we also use the quit channel. The user is supposed to send a message
// on the quit channel, and then close the connection
func ListenOnUDP(address string, readResults chan *UDPConnectionReadResult, quit chan struct{}) (*net.UDPConn, error) {
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
			select {
			case <-quit:
				return
			default:
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
		}
	}()

	return conn, nil
}
