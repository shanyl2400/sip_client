package model

type DeviceBrief struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Address  string `json:"address"`
}

type AddDeviceRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
}

type SendMessageRequest struct {
	MessageType string            `json:"message_type"`
	Values      map[string]string `json:"values"`
}
