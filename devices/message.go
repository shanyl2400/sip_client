package devices

type Message interface {
	Bytes() []byte
	ID() string
}
