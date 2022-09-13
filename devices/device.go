package devices

import (
	"sipsimclient/errors"
	"sipsimclient/model"
)

const (
	DeviceStateReady     DeviceState = "ready"
	DeviceStateConnected DeviceState = "connected"
	DeviceStateRegisting DeviceState = "registing"
	DeviceStateOnline    DeviceState = "online"
	DeviceStateOffline   DeviceState = "offline"

	DeviceStateErr      DeviceState = "error"
	DeviceStateUnauthed DeviceState = "unauthed"

	NetProtocolTCP = "tcp"
	NetProtocolUDP = "udp"
)

type DeviceState string
type NetProtocol string

type Device interface {
	Connect() error
	Disconnect() error
	Send(msg Message) error

	Logs(theme model.Theme) ([]string, error)

	Name() string
	Address() string
	Protocol() NetProtocol
	State() DeviceState
}

func NewDevice(req AddDeviceRequest) (Device, error) {
	if req.Protocol == NetProtocolTCP || req.Protocol == NetProtocolUDP {
		return createSocketDevice(req)
	}
	return nil, errors.ErrUnsupportedProtocol
}
