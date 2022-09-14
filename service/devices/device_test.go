package devices

import (
	"sipsimclient/config"
	"sipsimclient/devices/message"
	"testing"

	"github.com/jart/gosip/sip"
)

func TestConnectDevice(t *testing.T) {
	config.Set(&config.Config{
		ServerSocketHost: "127.0.0.1",
		ServerSocketPort: 5061,
	})

	device, err := NewDevice(AddDeviceRequest{
		Name:     "device1",
		Protocol: "tcp",
	})
	if err != nil {
		t.Fatal(err)
	}
	device.Connect()
	select {}
}

func TestSIPMessage(t *testing.T) {
	config.Set(&config.Config{
		ServerSocketHost: "127.0.0.1",
		ServerSocketPort: 8051,
		SIPScope:         "123123123",
	})
	heartBeat := message.NewHeartBeatMessage("device1", "127.0.0.1", 67535)
	t.Log(string(heartBeat.Bytes()))

	msg, err := sip.ParseMsg(heartBeat.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msg)
}
