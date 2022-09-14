package message

import (
	"bytes"

	"github.com/jart/gosip/sip"
)

type cancelMessage struct {
	host string
	port int
	name string
	id   string
}

func (h *cancelMessage) ID() string {
	return h.id
}

func (h *cancelMessage) Bytes() []byte {
	msg := defaultMessage(messageValue{
		id:     h.id,
		host:   h.host,
		port:   h.port,
		name:   h.name,
		method: sip.MethodCancel,
	})
	var b bytes.Buffer
	msg.Append(&b)
	return b.Bytes()
}
func NewCancelMessage(name, host string, port int) *cancelMessage {
	return &cancelMessage{
		host: host,
		port: port,
		name: name,
		id:   generateCallID(),
	}
}
