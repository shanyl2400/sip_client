package message

import (
	"sipsimclient/config"

	"github.com/jart/gosip/sip"
)

type messageValue struct {
	id   string
	host string
	port int
	name string

	method string
}

func defaultMessage(h messageValue) *sip.Msg {
	sipurl := &sip.URI{
		Scheme: "sip",
		Host:   config.Get().ServerSocketHost,
		Port:   uint16(config.Get().ServerSocketPort),
	}
	return &sip.Msg{
		Method:     h.method,
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
		MaxForwards: 70,
		UserAgent:   "CAROT-SIP",
	}
}
