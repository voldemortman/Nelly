package nelly

type Bridge struct {
	deviceToDeviceCommunicatorMap map[string]*DeviceCommunicator
	deviceDestructor              DeviceDestructor
	isStarted                     bool
}

func (bridge *Bridge) StartBridge() {
	for device := range bridge.deviceToDeviceCommunicatorMap {
		bridge.startDevice(device)
	}
	bridge.isStarted = true
}

func (bridge *Bridge) StopBridge() error {
	for device := range bridge.deviceToDeviceCommunicatorMap {
		if err := bridge.RemoveDevice(device); err != nil {
			sugar.Errorf("failed to stop device %s", device)
			return err
		}
	}
	return nil
}

func (bridge *Bridge) AddDevice(deviceName string) error {
	communicator, err := BuildDeviceHandleCommunication(deviceName, nil)
	if err != nil {
		sugar.Errorf("failed to build device communicator on device %s", deviceName)
		return err
	}
	bridge.deviceToDeviceCommunicatorMap[deviceName] = communicator
	if bridge.isStarted {
		bridge.startDevice(deviceName)
	}
	return nil
}

func (bridge *Bridge) RemoveDevice(deviceName string) error {
	delete(bridge.deviceToDeviceCommunicatorMap, deviceName)
	return bridge.deviceDestructor.DestructDevice(deviceName)
}

func (bridge *Bridge) startDevice(deviceName string) {
	communicator := bridge.deviceToDeviceCommunicatorMap[deviceName]
	writerQuit := make(chan bool)
	sourceQuit := make(chan bool)

	bridge.startWritingOnDevice(deviceName, writerQuit)
	sourceListenerQuit := bridge.startListeningOnDevice(deviceName, sourceQuit)

	bridge.deviceDestructor.DeviceToQuitChannelMap[deviceName] = []chan bool{sourceQuit, sourceListenerQuit, writerQuit}
	bridge.deviceDestructor.DeviceToHandleMap[deviceName] = communicator.Handle
}

func (bridge *Bridge) startWritingOnDevice(deviceName string, writerQuit chan bool) {
	communicator := bridge.deviceToDeviceCommunicatorMap[deviceName]
	// TODO: code smell, this returns a channel that we don't use
	communicator.PacketWriter.Start(writerQuit)
}

func (bridge *Bridge) startListeningOnDevice(deviceName string, sourceQuit chan bool) chan bool {
	communicator := bridge.deviceToDeviceCommunicatorMap[deviceName]
	packetStream := communicator.PacketSniffer.Start(sourceQuit)
	sourceListenerQuit := make(chan bool)
	go func() {
		wasQuitCalled := false
		for !wasQuitCalled {
			select {
			case packet := <-packetStream:
				// TODO: maybe this part should be a filter, maybe device sniffer should hold a referance to the future appropriate device retriever struct
				// This is a POC. It is hard coded to being a hub.
				for device, deviceCommunicator := range bridge.deviceToDeviceCommunicatorMap {
					if device != deviceName {
						deviceCommunicator.PacketWriter.source <- packet
					}
				}
			case <-sourceListenerQuit:
				sugar.Infof("Stopped listening and forwarding on device %s", deviceName)
				wasQuitCalled = true
			}
		}
	}()
	return sourceListenerQuit
}
