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
	msg := defaultMessage(messageValue{
		id:     h.id,
		host:   h.host,
		port:   h.port,
		name:   h.name,
		method: sip.MethodRegister,
	})
	msg.Authorization = h.auth()
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
