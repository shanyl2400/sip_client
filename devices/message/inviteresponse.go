package message

import (
	"bytes"
	"sipsimclient/config"

	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
)

// type SDP struct {
// 	Origin   Origin      // This must always be present
// 	Addr     string      // Connect to this IP; never blank (from c=)
// 	Audio    *Media      // Non-nil if we can establish audio
// 	Video    *Media      // Non-nil if we can establish video
// 	Session  string      // s= Session Name (default "-")
// 	Time     string      // t= Active Time (default "0 0")
// 	Ptime    int         // Transmit frame every N milliseconds (default 20)
// 	SendOnly bool        // True if 'a=sendonly' was specified in SDP
// 	RecvOnly bool        // True if 'a=recvonly' was specified in SDP
// 	Attrs    [][2]string // a= lines we don't recognize
// 	Other    [][2]string // Other description
// }

type inviteResponse struct {
	ip       string
	port     int
	protocol string
	msg      *sip.Msg
}

func (h *inviteResponse) ID() string {
	return h.msg.CallID
}

func (h *inviteResponse) Bytes() []byte {
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
		Payload:     h.mockPayload(),
	}
	var b bytes.Buffer
	msg.Append(&b)
	return b.Bytes()
}
func (h *inviteResponse) mockPayload() *sdp.SDP {
	// var other [][2]string
	// other = append(other, [2]string{"y", fmt.Sprintf("%010d", ssrc)})
	/*判断是tcp就请求tcp流*/
	var proto string
	var attr [][2]string
	if h.protocol == "tcp" {
		proto = "TCP/RTP/AVP"
		attr = append(attr, [2]string{"setup", "passive"}) //--tcp传输时有 active表示发送者是客户端 passive表示发送者是服务端
		attr = append(attr, [2]string{"connection", "new"})
	} else {
		proto = "RTP/AVP"
	}
	return &sdp.SDP{
		Origin: sdp.Origin{
			User:    config.Get().SIPScope,
			ID:      "0",
			Version: "0",
			Addr:    h.ip,
		},
		Addr:    h.ip,
		Session: "Play",
		Video: &sdp.Media{
			Proto: proto,
			Port:  uint16(h.port),
			Codecs: []sdp.Codec{{
				PT:   96,
				Name: "PS",
				Rate: 90000,
			}},
		},
		RecvOnly: true,
		Attrs:    attr,
		// Other:    other,
	}
}
func NewInviteResponse(msg *sip.Msg, ip string, port int, protocol string) *inviteResponse {
	return &inviteResponse{
		ip:       ip,
		port:     port,
		msg:      msg,
		protocol: protocol,
	}
}
