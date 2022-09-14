package devices

import "sipsimclient/service/devices/message"

const (
	MessageHeartBeat MessageType = "heartbeeat"
	MessageRegister  MessageType = "register"
	MessageBye       MessageType = "bye"
	MessageCancel    MessageType = "cancel"
	MessageRaw       MessageType = "raw"
)

type MessageType string

func createMessage(device Device, messageType MessageType, val map[string]string) Message {
	switch messageType {
	case MessageHeartBeat:
		return message.NewHeartBeatMessage(device.Name(), device.Host(), device.Port())
	case MessageRegister:
		return message.NewRegisterMessage(device.Name(), val["password"], device.Host(), device.Port())
	case MessageBye:
		return message.NewByeMessage(device.Name(), device.Host(), device.Port())
	case MessageCancel:
		return message.NewCancelMessage(device.Name(), device.Host(), device.Port())
	case MessageRaw:
		return message.NewRawMessage(val["id"], []byte(val["data"]))
	}
	return nil
}
