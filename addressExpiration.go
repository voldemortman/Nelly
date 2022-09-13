package nelly

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/gopacket"
)

type AddressExpiration struct {
	endpointToTimeStampMap map[gopacket.Endpoint]time.Time
	limit                  time.Duration
	endpointType           gopacket.EndpointType
}

func (expirations *AddressExpiration) UpdateAddressTimeStamp(address gopacket.Endpoint) error {
	err := validateEndpointType(expirations.endpointType, address)
	if err != nil {
		return err
	}
	expirations.endpointToTimeStampMap[address] = time.Now()
	return nil
}

func (expirations *AddressExpiration) TryGetAddressTimeStamp(address gopacket.Endpoint, timeStamp *time.Time) (bool, error) {
	err := validateEndpointType(expirations.endpointType, address)
	if err != nil {
		return false, err
	}
	now := time.Now()
	if addressTimeStamp, ok := expirations.endpointToTimeStampMap[address]; ok && now.Sub(addressTimeStamp) >= expirations.limit {
		*timeStamp = addressTimeStamp
		return true, nil
	}
	return false, nil
}

func validateEndpointType(endpointType gopacket.EndpointType, endpoint gopacket.Endpoint) error {
	if endpoint.EndpointType() == endpointType {
		return nil
	}
	return errors.New(fmt.Sprintf("recieved endpoint is not of same type as defined in struct. Recieved type: %s. Wanted type %s",
		endpoint.EndpointType(), endpointType))
}
