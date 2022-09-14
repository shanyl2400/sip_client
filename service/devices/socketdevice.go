package devices

import (
	"fmt"
	"io"
	"net"
	"sipsimclient/config"
	"sipsimclient/log"
	"sipsimclient/model"
	"sipsimclient/service/devices/message"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jart/gosip/sip"
)

type socketDevice struct {
	name     string
	addr     string
	password string

	state  DeviceState
	logger *DeviceLogger

	protocol NetProtocol

	conn net.Conn

	host string
	port int

	pendingResponse      map[string]Message
	pendingResponseMutex sync.Mutex

	ticker    *time.Ticker
	quit      chan struct{}
	socketMsg chan string
}

func (td *socketDevice) Connect() error {
	if td.state != DeviceStateReady && td.state != DeviceStateOffline &&
		td.state != DeviceStateErr {
		log.Infof("device %v is already connected", td.name)
		return nil
	}
	td.state = DeviceStateConnected
	//TODO: implement it
	addr := fmt.Sprintf("%v:%v", config.Get().ServerSocketHost, config.Get().ServerSocketPort)
	conn, err := net.DialTimeout(string(td.Protocol()), addr, 5*time.Second) //创建套接字,连接服务器,设置超时时间
	if err != nil {
		return err
	}
	td.addr = conn.LocalAddr().String()
	td.conn = conn
	td.setHostPort(td.addr)

	go td.receive()
	go td.handleMessage()

	//register
	err = td.registerDevice()
	if err != nil {
		td.logger.Errorf("register device failed, err: %v", err)
		td.state = DeviceStateErr
		return err
	}
	td.state = DeviceStateRegisting

	return nil
}
func (td *socketDevice) Disconnect() error {
	if td.state == DeviceStateOffline {
		log.Infof("device %v is already disconnected", td.name)
		return nil
	}
	//TODO: implement it
	log.Debugf("%v prepare to close", td.name)
	td.conn.Close()
	log.Debugf("%v connect closed", td.name)
	td.quit <- struct{}{}
	log.Debugf("%v quit handler", td.name)
	td.ticker.Stop()
	log.Debugf("%v stop ticker", td.name)
	td.state = DeviceStateOffline
	return nil
}

func (td *socketDevice) Send(msg Message) error {
	if td.state != DeviceStateConnected && td.state != DeviceStateRegisting && td.state != DeviceStateOnline && td.state != DeviceStateUnauthed {
		log.Infof("device %v is not online", td.name)
		return nil
	}
	msgBytes := msg.Bytes()
	td.logger.Send("send message", string(msgBytes))
	_, err := td.conn.Write(msgBytes)
	return err
}

func (td *socketDevice) SendForResponse(msg Message) error {
	err := td.Send(msg)
	if err != nil {
		return err
	}
	td.pendingResponseMutex.Lock()
	defer td.pendingResponseMutex.Unlock()
	td.pendingResponse[msg.ID()] = msg
	return nil
}

func (td *socketDevice) Logs(theme model.Theme, start, end time.Time) ([]*model.DeviceLog, error) {
	//TODO: implement it
	return td.logger.RangeLogs(theme, start, end)
}

func (td *socketDevice) Update(password string, protocol NetProtocol) error {
	if password == "" && protocol == "" {
		log.Debugf("%v nothing to update", td.name)
		return nil
	}
	if td.state != DeviceStateErr && td.state != DeviceStateOffline {
		err := td.Disconnect()
		if err != nil {
			log.Warnf("%v disconnect failed, err: %v", td.name, err)
			return err
		}
	}
	if password != "" {
		td.password = password
	}
	if protocol != "" {
		td.protocol = protocol
	}

	err := td.Connect()
	if err != nil {
		log.Warnf("%v reconnect failed, err: %v", td.name, err)
		return err
	}

	return nil
}

func (td *socketDevice) Name() string {
	return td.name
}
func (td *socketDevice) Address() string {
	return td.addr
}
func (td *socketDevice) Protocol() NetProtocol {
	return td.protocol
}
func (td *socketDevice) State() DeviceState {
	return td.state
}
func (td *socketDevice) Host() string {
	return td.host
}
func (td *socketDevice) Port() int {
	return td.port
}
func (td *socketDevice) setHostPort(addr string) {
	addrParts := strings.Split(addr, ":")
	if len(addrParts) < 2 {
		td.host = addr
		return
	}
	td.host = addrParts[0]
	td.port, _ = strconv.Atoi(addrParts[1])
}

func (td *socketDevice) registerDevice() error {
	msg := message.NewRegisterMessage(td.name, td.password, td.host, td.port)
	err := td.SendForResponse(msg)
	if err != nil {
		return err
	}
	return nil
}

func (td *socketDevice) onReceive(msg *sip.Msg, raw string) {
	//TODO: implement it
	if msg.Request == nil {
		//response
		td.handleResponse(msg, raw)
	} else {
		//request
		td.handleRequest(msg, raw)
	}
}

func (td *socketDevice) handleResponse(msg *sip.Msg, raw string) {
	td.pendingResponseMutex.Lock()
	_, ok := td.pendingResponse[msg.CallID]
	delete(td.pendingResponse, msg.CallID)
	td.pendingResponseMutex.Unlock()

	if !ok {
		//无效响应
		// msgJSON, _ := json.Marshal(msg)
		td.logger.Receive("Unknown response", raw)
		return
	}

	switch msg.CSeqMethod {
	case sip.MethodRegister:
		//响应Register, 修改device状态
		if msg.Status == sip.StatusOK {
			log.Debugf("%v: receive ok response", td.name)
			td.state = DeviceStateOnline
		} else if msg.Status == sip.StatusUnauthorized {
			log.Debugf("%v: receive unauthorized response", td.name)
			td.state = DeviceStateUnauthed
		} else {
			log.Debugf("%v: receive error response", td.name)
			td.state = DeviceStateErr
		}
	default:
		td.logger.Receive("Unhandled response", raw)
	}
}
func (td *socketDevice) handleRequest(msg *sip.Msg, raw string) {
	//处理请求
	switch msg.Method {
	case sip.MethodInvite:
		log.Debugf("%v: receive invite", td.name)
		td.SendForResponse(message.NewInviteResponse(msg, td.host, td.port, string(td.protocol)))
	case sip.MethodBye:
		log.Debugf("%v: receive bye", td.name)
		td.SendForResponse(message.NewOKResponseMessage(msg))
	case sip.MethodCancel:
		log.Debugf("%v: receive cancel", td.name)
		td.SendForResponse(message.NewOKResponseMessage(msg))
	default:
		td.logger.Receive("Unhandled request", raw)
	}
}

func (td *socketDevice) handleMessage() {
	for {
		select {
		case text := <-td.socketMsg:
			msg, err := td.parseSIPMsg(text)
			if err != nil {
				td.logger.Receivef("Can't parse sip message, err: %v", text, err)
				td.state = DeviceStateErr
				log.Warnf("parse message failed, err: %v", err)
				return
			}
			td.onReceive(msg, text)
		case <-td.ticker.C:
			//send liveness message
			err := td.Send(message.NewHeartBeatMessage(td.name, td.host, td.port))
			if err != nil {
				log.Warnf("send failed: %v", err)
			}
		case <-td.quit:
			return
		}
	}
}
func (td *socketDevice) parseSIPMsg(text string) (*sip.Msg, error) {
	buf := []byte(text)
	msg, err := sip.ParseMsg(buf)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (td *socketDevice) receive() {
	var buf []byte = make([]byte, 4096)
	for {
		n, err := td.conn.Read(buf) //将客户端发送的数据写入buf中
		if err != nil {
			if err == io.EOF {
				td.logger.Info("socket closed")
				td.Disconnect()
				break
			} else {
				td.logger.Errorf("read socket message failed, err: %v", err)
				td.state = DeviceStateErr
				break
			}
		}
		msg := string(buf[:n])
		td.logger.Receive("receive a message", msg)
		td.socketMsg <- string(msg)
	}
}

func createSocketDevice(req AddDeviceRequest) (*socketDevice, error) {
	logger, err := NewLogger(req.Name)
	if err != nil {
		return nil, err
	}

	//TODO: implement it
	return &socketDevice{
		name:            req.Name,
		password:        req.Password,
		state:           DeviceStateReady,
		protocol:        req.Protocol,
		logger:          logger,
		ticker:          time.NewTicker(30 * time.Second),
		quit:            make(chan struct{}),
		socketMsg:       make(chan string, 1),
		pendingResponse: make(map[string]Message),
	}, nil
}
