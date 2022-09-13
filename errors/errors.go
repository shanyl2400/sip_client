package errors

import "errors"

var (
	ErrDuplicateDeviceName = errors.New("duplicate device name")
	ErrDeviceNotExists     = errors.New("devices is not exists")

	ErrUnsupportedProtocol = errors.New("unsupported protocol")

	ErrHttpRequestErr = errors.New("http error")

	ErrUserNotExists      = errors.New("user is not exists")
	ErrInvalidMessageType = errors.New("invalid message type")
)
