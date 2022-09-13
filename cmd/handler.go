package cmd

import (
	"encoding/json"
	"fmt"
	"sipsimclient/devices"
	"sipsimclient/http"
	"sipsimclient/model"
	"strconv"
	"strings"
	"time"
)

func (p *Prompt) handleDevice(args []string) {
	if len(args) < 1 {
		fmt.Println("devices need args")
		return
	}
	switch args[0] {
	case "list":
		p.listDevices()
	case "add":
		p.addDevice(args[1:])
	case "connect":
		p.connect(args[1:])
	case "disconnect":
		p.disconnect(args[1:])
	case "remove":
		p.remove(args[1:])
	case "send":
		p.send(args[1:])
	case "logs":
		p.logs(args[1:])
	default:
		fmt.Println("unknown command")
	}
}

func (p *Prompt) handleHTTP(args []string) {
	if len(args) < 2 {
		fmt.Println("http need args")
		return
	}
	switch args[0] {
	case "session":
		p.handleSessionCommand(args[1:])
	case "device":
		p.handleDeviceCommand(args[1:])
	case "ptzv":
		p.ptzv(args[1:])
	}
}

func (p *Prompt) handleSessionCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("session args empty")
		return
	}
	switch args[0] {
	case "create":
		p.createSession(args[1:])
	case "release":
		p.releaseSession(args[1:])
	case "recover":
		p.recoverSession(args[1:])
	default:
		fmt.Println("unknown command")
	}
}

func (p *Prompt) handleDeviceCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("session args empty")
		return
	}
	switch args[0] {
	case "add":
		p.registerDevice(args[1:])
	case "get":
		p.getDeviceStat(args[1:])
	default:
		fmt.Println("unknown command")
	}
}

func (p *Prompt) createSession(args []string) {
	req, err := p.parseCreateSession(args)
	if err != nil {
		fmt.Println("invalid create session params")
		return
	}
	session, err := p.sipClient.CreateSession(req)
	if err != nil {
		fmt.Println("create session failed, err:", err)
		return
	}
	sessionJSON, _ := json.MarshalIndent(session, "", "\t")
	fmt.Println("response:\n", string(sessionJSON))
}
func (p *Prompt) releaseSession(args []string) {
	req, err := p.parseReleaseSession(args)
	if err != nil {
		fmt.Println("invalid release session params")
		return
	}
	err = p.sipClient.ReleaseSession(req)
	if err != nil {
		fmt.Println("release session failed, err:", err)
		return
	}
	fmt.Println("release session successfully")
}
func (p *Prompt) recoverSession(args []string) {
	req, err := p.parseRecoverSession(args)
	if err != nil {
		fmt.Println("invalid recover session params")
		return
	}
	err = p.sipClient.RecoverSession(req)
	if err != nil {
		fmt.Println("recover session failed, err:", err)
		return
	}
	fmt.Println("recover session successfully")
}

func (p *Prompt) registerDevice(args []string) {
	req := p.parseAddDevice(args)
	err := p.sipClient.AddDevice(req)

	if err != nil {
		fmt.Println("add device failed, err:", err)
		return
	}
	fmt.Println("add device successfully")
}
func (p *Prompt) getDeviceStat(args []string) {
	argsMap := p.parseArgs(args)
	deviceID := argsMap["d"]
	device, err := p.sipClient.GetDeviceStat(deviceID)
	if err != nil {
		fmt.Println("ptzv failed, err:", err)
		return
	}
	sessionJSON, _ := json.MarshalIndent(device, "", "\t")
	fmt.Println("device:\n", sessionJSON)
}

func (p *Prompt) ptzv(args []string) {
	req, err := p.parseCameraCtrl(args)
	if err != nil {
		fmt.Println("invalid ptzv params")
		return
	}
	err = p.sipClient.PTZV(req)
	if err != nil {
		fmt.Println("ptzv failed, err:", err)
		return
	}
	fmt.Println("ptzv successfully")
}

func (p *Prompt) listDevices() {
	devices := p.deviceManager.List()
	fmt.Println(devices.String())
}
func (p *Prompt) addDevice(args []string) {
	params := p.parseArgs(args)
	name := params["n"]
	protocl := params["p"]
	password := params["w"]
	err := p.deviceManager.Add(devices.AddDeviceRequest{
		Name:     name,
		Password: password,
		Protocol: devices.NetProtocol(protocl),
	})
	if err != nil {
		fmt.Println("Add device failed, err:", err)
		return
	}
	fmt.Println("add device successfully")
}

func (p *Prompt) connect(args []string) {
	if len(args) != 1 {
		fmt.Println("connect args only need one")
		return
	}
	err := p.deviceManager.Connect(args[0])
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}
	fmt.Println("connect successfully")
}

func (p *Prompt) disconnect(args []string) {
	if len(args) != 1 {
		fmt.Println("disconnect args only need one")
		return
	}
	err := p.deviceManager.Disconnect(args[0])
	if err != nil {
		fmt.Println("disconnect failed, err:", err)
		return
	}
	fmt.Println("disconnect successfully")
}

func (p *Prompt) remove(args []string) {
	if len(args) != 1 {
		fmt.Println("remove args only need one")
		return
	}
	err := p.deviceManager.Remove(args[0])
	if err != nil {
		fmt.Println("remove failed, err:", err)
		return
	}
	fmt.Println("remove successfully")
}

func (p *Prompt) send(args []string) {
	fmt.Println("This function has not implement")
}

func (p *Prompt) logs(args []string) {
	if len(args) != 2 {
		fmt.Println("remove args only need one")
		return
	}
	end := time.Now()
	start := end.Add(-time.Minute * 5)
	logs, err := p.deviceManager.Logs(args[0], model.Theme(args[1]), start, end)
	if err != nil {
		fmt.Println("remove failed, err:", err)
		return
	}
	for _, log := range logs {
		fmt.Println(log.String())
	}
}

func (p *Prompt) parseArgs(args []string) map[string]string {
	key := ""
	ans := make(map[string]string)
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			key = args[i][1:]
			ans[key] = ""
			continue
		}
		if key != "" {
			ans[key] = args[i]
			key = ""
		}
	}
	return ans
}

func (p *Prompt) parseCreateSession(args []string) (*http.CreateSessionRequest, error) {
	params := p.parseArgs(args)
	deviceID := params["d"]
	channelID := params["c"]
	ip := params["i"]
	portStr := params["p"]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	return &http.CreateSessionRequest{
		DeviceID:  deviceID,
		ChannelID: channelID,
		IP:        ip,
		UdpPort:   port,
		TcpPort:   port,
	}, nil
}

func (p *Prompt) parseReleaseSession(args []string) (*http.ReleaseSessionRequest, error) {
	params := p.parseArgs(args)
	deviceID := params["d"]
	dialogID := params["di"]
	channelID := params["c"]
	ssrcStr := params["s"]
	ip := params["i"]
	portStr := params["p"]
	sdp := params["sdp"]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	ssrc, err := strconv.Atoi(ssrcStr)
	if err != nil {
		return nil, err
	}
	return &http.ReleaseSessionRequest{
		DeviceID:  deviceID,
		DialogID:  dialogID,
		ChannelID: channelID,
		Ssrc:      uint32(ssrc),
		DstIP:     ip,
		DstPort:   port,
		SDP:       sdp,
	}, nil
}

func (p *Prompt) parseRecoverSession(args []string) (*http.RecoverSessionRequest, error) {
	params := p.parseArgs(args)
	deviceID := params["d"]
	dialogID := params["di"]
	channelID := params["c"]
	ssrcStr := params["s"]
	ip := params["i"]
	portStr := params["p"]
	sdp := params["sdp"]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	ssrc, err := strconv.Atoi(ssrcStr)
	if err != nil {
		return nil, err
	}
	return &http.RecoverSessionRequest{
		DeviceID:  deviceID,
		DialogID:  dialogID,
		ChannelID: channelID,
		Ssrc:      uint32(ssrc),
		DstIP:     ip,
		DstPort:   port,
		SDP:       sdp,
	}, nil
}

func (p *Prompt) parseCameraCtrl(args []string) (*http.CameraCtrlRequest, error) {
	params := p.parseArgs(args)
	cameraId := params["c"]
	actionStr := params["a"]
	paramStr := params["p"]
	stepStr := params["s"]

	action, err := strconv.Atoi(actionStr)
	if err != nil {
		return nil, err
	}
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return nil, err
	}
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		return nil, err
	}
	return &http.CameraCtrlRequest{
		StrCamID: cameraId,
		DwAction: action,
		DwParam:  param,
		DwStep:   step,
	}, nil
}

func (p *Prompt) parseAddDevice(args []string) *http.AddDeviceRequest {
	params := p.parseArgs(args)
	deviceID := params["d"]
	password := params["p"]
	scope := params["s"]

	return &http.AddDeviceRequest{
		DeviceID: deviceID,
		Password: password,
		Scope:    scope,
	}
}
