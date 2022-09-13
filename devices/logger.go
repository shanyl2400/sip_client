package devices

import (
	"fmt"
	"sipsimclient/model"
	"sipsimclient/repository"
)

type DeviceLogger struct {
	deviceName string
}

func (d *DeviceLogger) Info(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeInfo,
		Message:    text,
	})
}
func (d *DeviceLogger) Infof(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeInfo,
		Message:    fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Warn(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeWarn,
		Message:    text,
	})
}
func (d *DeviceLogger) Warnf(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeWarn,
		Message:    fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Error(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeError,
		Message:    text,
	})
}
func (d *DeviceLogger) Errorf(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeError,
		Message:    fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Send(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeSend,
		Message:    text,
	})
}
func (d *DeviceLogger) Sendf(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeSend,
		Message:    fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Receive(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeRecevice,
		Message:    text,
	})
}
func (d *DeviceLogger) Receivef(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeRecevice,
		Message:    fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Logs(theme model.Theme) ([]string, error) {
	logs, err := repository.GetDeviceLogRepository().Query(d.deviceName, theme)
	if err != nil {
		fmt.Println("query logs failed", err)
		return nil, err
	}
	out := make([]string, len(logs))
	for i := range logs {
		out[i] = fmt.Sprintf("%v device: %v, theme: %v, msg: %v\n",
			logs[i].CreatedAt.Format("2006/01/02 03:04:05.9999"),
			logs[i].DeviceName, logs[i].Theme, logs[i].Message)
	}
	return out, nil
}

func NewLogger(deviceName string) (*DeviceLogger, error) {
	return &DeviceLogger{
		deviceName: deviceName,
	}, nil
}

func ReleaseLogger(deviceName string) {
	repository.GetDeviceLogRepository().DeleteAll(deviceName)
}
