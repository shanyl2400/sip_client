package sip

import (
	"fmt"
	"net/http"
	"sipsimclient/config"
	"sipsimclient/errors"
	"sipsimclient/log"

	"github.com/go-resty/resty/v2"
)

type CreateSessionRequest struct {
	DeviceID  string `json:"deviceid"`
	ChannelID string `json:"channelid"`
	IP        string `json:"ip"`
	UdpPort   int    `json:"udpport"`
	TcpPort   int    `json:"tcpport"`
}

type ReleaseSessionRequest struct {
	DeviceID  string `json:"deviceid"`
	DialogID  string `json:"dialogid"`
	ChannelID string `json:"channelid"`
	Ssrc      uint32 `json:"ssrc"`
	DstIP     string `json:"dstip"`
	DstPort   int    `json:"dstport"`
	SDP       string `json:"sdp"`
}

type RecoverSessionRequest struct {
	DeviceID  string `json:"deviceid"`
	DialogID  string `json:"dialogid"`
	ChannelID string `json:"channelid"`
	Ssrc      uint32 `json:"ssrc"`
	DstIP     string `json:"dstip"`
	DstPort   int    `json:"dstport"`
	SDP       string `json:"sdp"`
}
type AddDeviceRequest struct {
	DeviceID string `json:"device_id"`
	Password string `json:"passwd"`
	Scope    string `json:"scope"`
}

type CameraCtrlRequest struct {
	StrCamID string `json:"cameraId"`
	DwAction int    `json:"action"`
	DwParam  int    `json:"param"`
	DwStep   int    `json:"step"`
}

type Session struct {
	DeviceID  string `json:"deviceid"`
	DialogID  string `json:"dialogid"`
	ChannelID string `json:"channelid"`
	Ssrc      uint32 `json:"ssrc"`
	DstIP     string `json:"dstip"`
	DstPort   int    `json:"dstport"`
	SDP       string `json:"sdp"`
}

type SIPDeviceStat struct {
	SIPID      string `json:"sipid"`
	Transport  int    `json:"transport"`
	RemoteAddr string `json:"remoteaddr"`
	State      int    `json:"state"`
}

type SIPClient interface {
	CreateSession(req *CreateSessionRequest) (*Session, error)
	ReleaseSession(req *ReleaseSessionRequest) error
	RecoverSession(req *RecoverSessionRequest) error
	GetDeviceStat(deviceid string) (*SIPDeviceStat, error)
	AddDevice(req *AddDeviceRequest) error
	PTZV(req *CameraCtrlRequest) error
}

type SIPHttpClient struct {
	client *resty.Client
}

func (sc *SIPHttpClient) CreateSession(req *CreateSessionRequest) (*Session, error) {
	ans := new(Session)
	resp, err := sc.client.R().SetBody(req).
		SetResult(ans).Post("/api/v3/sip/createsession")
	if err != nil {
		return nil, err
	}
	if resp.RawResponse.StatusCode != http.StatusOK {
		log.Infof("request with failed, body:%v", string(resp.Body()))
		return nil, errors.ErrHttpRequestErr
	}

	return ans, nil
}
func (sc *SIPHttpClient) ReleaseSession(req *ReleaseSessionRequest) error {
	resp, err := sc.client.R().SetBody(req).Post("/api/v3/sip/releasesession")
	if err != nil {
		return err
	}
	if resp.RawResponse.StatusCode != http.StatusOK {
		log.Infof("request with failed, body:%v", string(resp.Body()))
		return errors.ErrHttpRequestErr
	}

	return nil
}
func (sc *SIPHttpClient) RecoverSession(req *RecoverSessionRequest) error {
	resp, err := sc.client.R().SetBody(req).Post("/api/v3/sip/recoversession")
	if err != nil {
		return err
	}
	if resp.RawResponse.StatusCode != http.StatusOK {
		log.Infof("request with failed, body:%v", string(resp.Body()))
		return errors.ErrHttpRequestErr
	}

	return nil
}
func (sc *SIPHttpClient) GetDeviceStat(deviceid string) (*SIPDeviceStat, error) {
	ans := new(SIPDeviceStat)
	resp, err := sc.client.R().SetQueryParam("deviceid", deviceid).
		SetResult(ans).Get("/api/v3/sip/devicestat")
	if err != nil {
		return nil, err
	}
	if resp.RawResponse.StatusCode != http.StatusOK {
		log.Infof("request with failed, body:%v", string(resp.Body()))
		return nil, errors.ErrHttpRequestErr
	}

	return ans, nil
}
func (sc *SIPHttpClient) AddDevice(req *AddDeviceRequest) error {
	resp, err := sc.client.R().SetFormData(map[string]string{
		"deviceId": req.DeviceID,
		"passwd":   req.Password,
		"scope":    req.Scope,
	}).Post("/api/device/permission")
	if err != nil {
		return err
	}
	if resp.RawResponse.StatusCode != http.StatusOK {
		log.Infof("request with failed, body:%v", string(resp.Body()))
		return errors.ErrHttpRequestErr
	}

	return nil
}
func (sc *SIPHttpClient) PTZV(req *CameraCtrlRequest) error {
	resp, err := sc.client.R().SetBody(req).Post("/api/v3/sip/ptz")
	if err != nil {
		return err
	}
	if resp.RawResponse.StatusCode != http.StatusOK {
		log.Infof("request with failed, body:%v", string(resp.Body()))
		return errors.ErrHttpRequestErr
	}

	return nil
}

func NewSIPClient() SIPClient {
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("%v:%v",
		config.Get().ServerHttpHost,
		config.Get().ServerHttpPort))
	return &SIPHttpClient{
		client: client,
	}
}
