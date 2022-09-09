package nelly

type Bridge struct {
	deviceToDeviceCommunicatorMap map[string]*DeviceCommunicator
	deviceDestrucotr              DeviceDestructor
	isStarted                     bool
}

func (bridge *Bridge) AddDevice(deviceName string) error {
	communicator, err := BuildDeviceCommunicator(deviceName, nil)
	if err != nil {
		sugar.Errorf("failed to build device communicator on device %s", deviceName)
		return err
	}
	bridge.deviceToDeviceCommunicatorMap[deviceName] = communicator
	if bridge.isStarted {
		// TODO
	}
	return nil
}

func (bridge *Bridge) RemoveDevice(deviceName string) error {
	delete(bridge.deviceToDeviceCommunicatorMap, deviceName)
	return bridge.deviceDestrucotr.DestructDevice(deviceName)
}
