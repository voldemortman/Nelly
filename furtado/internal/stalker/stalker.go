package stalker

import (
	"fmt"
	"furtado/pkg/netUtils"
)

type IStalker interface {
	StartStalking() error
	StopStalking()
}

type RemoteStalker struct {
	config RemoteStalkerConfig
	quit   chan struct{}
}

func (stalker *RemoteStalker) StartStalking() error {

	sourceAddress := fmt.Sprint(stalker.config.LocalIP, ":", stalker.config.Port)
	destAddress := fmt.Sprintf(stalker.config.BridgeIP, ":", stalker.config.Port)

	sendConn, err := netUtils.CreateUDPSocket(sourceAddress, destAddress)

	if err != nil {
		return err
	}

	device, err := netUtils.CreateNetworkDevice(stalker.config.Interface)
	if err != nil {
		return err
	}

	go func() {
		packets, err := device.ReadFromDevice()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			select {
			case <-stalker.quit:
				err := sendConn.Close()
				if err != nil {
					fmt.Println(err)
				}
				return
			case packet := <-packets:
				b, err := netUtils.SerializePacket(packet)
				if err != nil {
					fmt.Print(err)
				}
				sendConn.Write(b)
			}
		}
	}()

	readResults := make(chan *netUtils.UDPConnectionReadResult)
	quit := make(chan struct{})
	conn, err := netUtils.ListenOnUDP(sourceAddress, readResults, quit)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-stalker.quit:
				conn.Close()
				close(quit)
				return
			case data := <-readResults:
				device.SendToDevice(data.Buffer)
			}
		}
	}()

	return nil
}

func (stalker *RemoteStalker) StopStalking() {
	close(stalker.quit)
}
