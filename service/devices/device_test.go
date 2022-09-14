package devices

import (
	"fmt"
	"sipsimclient/config"
	"sipsimclient/service/devices/message"
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

func TestTrim(t *testing.T) {
	// text := "MESSAGE sip:d1@127.0.0.1 SIP/2.0\r\nFrom: <sip:34020000002000000001@3402000001>;tag=020505080200\r\nTo: <sip:d1@127.0.0.1>\r\nVia: SIP/2.0/UDP 127.0.0.1:5061;rport;branch=z9hG4bKe233e8731d17\r\nCall-ID: 302915660\r\nCSeq: 26901 MESSAGE\r\nMax-Forwards: 70\r\nContent-Type: Application/MANSCDP+xml\r\nContent-Length: 224\r\n\r\n  <Control>\n      <CmdType>DeviceControl</CmdType>\n      <SN>1</SN>\n      <DeviceID>d1</DeviceID>\n      <PTZCmd>A50F010800480005</PTZCmd>\n      <Info>\n          <ControlPriority>5</ControlPriority>\n      </Info>\n  </Control>\n"
	// text := "SIP/2.0 200 OK\r\nFrom: <sip:d5@127.0.0.1>;tag=060904060409\r\nTo: <sip:127.0.0.1:5061>;tag=48267f7b3479\r\nVia: SIP/2.0/TCP 127.0.0.1:63592;branch=z9hG4bK020807070902;branch=z9hG4bK020807070902\r\nCall-ID: 378295307\r\nCSeq: 29789 REGISTER\r\nUser-Agent: CAROT-SIP\r\nExpires: 0\r\nContent-Length: 0\r\n\r\n"
	// text := "MESSAGE sip:127.0.0.1:5061 SIP/2.0\r\nFrom: <sip:d5@127.0.0.1>;tag=040203020302\r\nTo: <sip:127.0.0.1:5061>\r\nVia: SIP/2.0/TCP 127.0.0.1:63592;branch=z9hG4bK030008010508\r\nContact: <sip:d5@127.0.0.1:63592>\r\nCall-ID: 966504778\r\nCSeq: 19211 REGISTER\r\nUser-Agent: CAROT-SIP\r\nMax-Forwards: 70\r\nExpires: 0\r\nContent-Length: 0\r\n\r\n "
	text := `Via: SIP/2.0/UDP 192.168.70.50:5060;rport;branch=z9hG4bK1444802676
From: <sip:34020000001320001001@3402000000>;tag=542557691
To: <sip:34020000002000000001@3402000000>
Call-ID: 323739489
CSeq: 20 MESSAGE
Content-Type: Application/MANSCDP+xml
Max-Forwards: 70
User-Agent: IP Camera
Content-Length:   164

<?xml version="1.0" encoding="UTF-8"?>
<Notify>
<CmdType>Keepalive</CmdType>
<SN>46591</SN>
<DeviceID>34020000001320001001</DeviceID>
<Status>OK</Status>
</Notify>
`
	td := new(socketDevice)

	fmt.Println(text)
	msg, err := td.parseSIPMsg(text)
	if err != nil {
		t.Error("parse message failed, ", err)
		return
	}
	t.Log(msg)
}
