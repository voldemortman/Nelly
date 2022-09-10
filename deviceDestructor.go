package nelly

import (
	"errors"
	"fmt"

	"github.com/google/gopacket/pcap"
)

// TODO: Think of way to make this composite design pattern, or maybe other patterns
type DeviceDestructor struct {
	DeviceToQuitChannelMap map[string][]chan bool
	DeviceToHandleMap      map[string]*pcap.Handle
}

func (destructor DeviceDestructor) DestructDevice(device string) error {
	if quitChannels, ok := destructor.DeviceToQuitChannelMap[device]; ok {
		for _, quit := range quitChannels {
			quitChan := quit
			go func() {
				quitChan <- true
			}()
		}
	} else {
		return errors.New(fmt.Sprintf("failed to get quit channels on device %s. Handle was not closed", device))
	}
	if handle, ok := destructor.DeviceToHandleMap[device]; ok {
		handle.Close()
	} else {
		return errors.New(fmt.Sprintf("failed to get handle on device %s. Quit channels were called", device))
	}
	delete(destructor.DeviceToHandleMap, device)
	delete(destructor.DeviceToQuitChannelMap, device)
	return nil
}
