package model

import (
	"fmt"
	"time"
)

type DeviceLog struct {
	DeviceName string    `json:"device_name"`
	Theme      Theme     `json:"theme"`
	Info       string    `json:"info"`
	Message    string    `json:"message,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (d *DeviceLog) String() string {
	return fmt.Sprintf("%v %v [%v] %v", d.CreatedAt.Format("2006/01/02 03:04:05.9999"), d.DeviceName, d.Theme, d.Message)
}
