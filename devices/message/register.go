package message

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sipsimclient/config"

	"github.com/jart/gosip/sip"
)

type registerMessage struct {
	host     string
	port     int
	name     string
	password string

	id string
}

func (h *registerMessage) ID() string {
	return h.id
}
func (h *registerMessage) Bytes() []byte {
	sipurl := &sip.URI{
		Scheme: "sip",
		Host:   config.Get().ServerSocketHost,
		Port:   uint16(config.Get().ServerSocketPort),
	}
	msg := sip.Msg{
		Method:     sip.MethodRegister,
		Request:    sipurl,
		CallID:     h.id,
		CSeq:       generateCSeq(),
		CSeqMethod: sip.MethodRegister,
		Via: &sip.Via{
			Transport: "TCP",
			Host:      h.host,
			Port:      uint16(h.port),
			Param: &sip.Param{
				Name:  "branch",
				Value: generateBranch(),
			}},
		From: &sip.Addr{
			Uri: &sip.URI{
				User: h.name,
				Host: h.host,
			},
			Param: &sip.Param{
				Name:  "tag",
				Value: generateTag(),
			}},
		To: &sip.Addr{
			Uri: sipurl,
		},
		Contact: &sip.Addr{
			Uri: &sip.URI{
				User: h.name,
				Host: h.host,
				Port: uint16(h.port),
			}},
		MaxForwards:   70,
		UserAgent:     "CAROT-SIP",
		Authorization: h.auth(),
	}
	var b bytes.Buffer
	msg.Append(&b)
	return b.Bytes()
}
func (h *registerMessage) auth() string {
	scope := config.Get().SIPScope
	nonce := fmt.Sprintf("%x", md5.Sum([]byte("this is carot SIP Serve")))
	md5UserRealmPwd := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", h.name, scope, h.password))))
	md5MethodURL := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", sip.MethodRegister, ""))))
	response := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", md5UserRealmPwd, nonce, md5MethodURL))))

	return fmt.Sprintf("Digest response=\"%v\",uri=\"\"", response)
}
func NewRegisterMessage(name, password, host string, port int) *registerMessage {
	return &registerMessage{
		host:     host,
		port:     port,
		name:     name,
		password: password,
		id:       generateCallID(),
	}
}
