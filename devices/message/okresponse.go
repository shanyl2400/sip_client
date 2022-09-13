package message

import (
	"bytes"

	"github.com/jart/gosip/sip"
)

type okResponseMessage struct {
	msg *sip.Msg
}

func (h *okResponseMessage) ID() string {
	return h.msg.CallID
}

func (h *okResponseMessage) Bytes() []byte {
	msg := sip.Msg{
		Status:      sip.StatusOK,
		CallID:      h.msg.CallID,
		CSeq:        h.msg.CSeq,
		CSeqMethod:  h.msg.Method,
		Via:         h.msg.Via,
		From:        h.msg.From,
		To:          h.msg.To,
		Contact:     h.msg.Contact,
		MaxForwards: 70,
		UserAgent:   "CAROT-SIP",
	}
	var b bytes.Buffer
	msg.Append(&b)
	return b.Bytes()
}
func NewOKResponseMessage(msg *sip.Msg) *okResponseMessage {
	return &okResponseMessage{
		msg: msg,
	}
}
