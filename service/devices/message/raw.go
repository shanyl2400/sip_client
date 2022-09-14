package message

type rawMessage struct {
	callID string
	raw    []byte
}

func (h *rawMessage) ID() string {
	return h.callID
}

func (h *rawMessage) Bytes() []byte {
	return h.raw
}
func NewRawMessage(id string, raw []byte) *rawMessage {
	return &rawMessage{
		callID: id,
		raw:    raw,
	}
}
