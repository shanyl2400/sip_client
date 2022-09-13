package message

import (
	"bytes"

	"github.com/jart/gosip/sip"
)

type byeMessage struct {
	host string
	port int
	name string
	id   string
}

func (h *byeMessage) ID() string {
	return h.id
}

func (h *byeMessage) Bytes() []byte {
	msg := defaultMessage(messageValue{
		id:     h.id,
		host:   h.host,
		port:   h.port,
		name:   h.name,
		method: sip.MethodBye,
	})
	var b bytes.Buffer
	msg.Append(&b)
	return b.Bytes()
}
func NewByeMessage(name, host string, port int) *byeMessage {
	return &byeMessage{
		host: host,
		port: port,
		name: name,
		id:   generateCallID(),
	}
}
