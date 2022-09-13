package message

import (
	"bytes"
	"encoding/hex"
	"math/rand"
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
	msg := defaultMessage(messageValue{
		id:     h.id,
		host:   h.host,
		port:   h.port,
		name:   h.name,
		method: sip.MethodMessage,
	})
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
