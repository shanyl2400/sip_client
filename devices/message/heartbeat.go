package message

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"sipsimclient/config"
	"strconv"

	"github.com/jart/gosip/sip"
)

type heartBeatMessage struct {
	host string
	port int
	name string
	id   string
}

func (h *heartBeatMessage) ID() string {
	return h.id
}

func (h *heartBeatMessage) Bytes() []byte {

	sipurl := &sip.URI{
		Scheme: "sip",
		Host:   config.Get().ServerSocketHost,
		Port:   uint16(config.Get().ServerSocketPort),
	}

	msg := sip.Msg{
		Method:     sip.MethodMessage,
		Request:    sipurl,
		CallID:     h.id,
		CSeq:       generateCSeq(),
		CSeqMethod: sip.MethodNotify,
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
		MaxForwards: 70,
		UserAgent:   "CAROT-SIP",
	}
	var b bytes.Buffer
	msg.Append(&b)
	return b.Bytes()
}
func NewHeartBeatMessage(name, host string, port int) *heartBeatMessage {
	return &heartBeatMessage{
		host: host,
		port: port,
		name: name,
		id:   generateCallID(),
	}
}

func generateBranch() string {
	return "z9hG4bK" + generateTag()
}

func generateCallID() string {
	lol := randomBytes(9)
	var str string
	for _, v := range lol {
		str += strconv.Itoa(int(v))
	}
	return str
}
func generateTag() string {
	return hex.EncodeToString(randomBytes(6))
}
func randomBytes(l int) (b []byte) {
	b = make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte(rand.Intn(10))
	}
	return
}

// Generate a secure random number between 0 and 50,000.
func generateCSeq() int {
	return rand.Int() % 50000
}
