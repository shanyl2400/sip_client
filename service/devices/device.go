package devices

import (
	"sipsimclient/errors"
	"sipsimclient/model"
	"time"
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

	Logs(theme model.Theme, start, end time.Time) ([]*model.DeviceLog, error)

	Name() string
	Address() string
	Protocol() NetProtocol
	State() DeviceState

	Host() string
	Port() int

	Update(password string, protocol NetProtocol) error
}

func NewDevice(req AddDeviceRequest) (Device, error) {
	if req.Protocol == NetProtocolTCP || req.Protocol == NetProtocolUDP {
		return createSocketDevice(req)
	}
	return nil, errors.ErrUnsupportedProtocol
}
