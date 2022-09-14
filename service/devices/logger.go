package devices

import (
	"fmt"
	"sipsimclient/model"
	"sipsimclient/repository"
	"time"
)

type DeviceLogger struct {
	deviceName string
}

func (d *DeviceLogger) Info(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeInfo,
		Info:       text,
	})
}
func (d *DeviceLogger) Infof(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeInfo,
		Info:       fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Warn(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeWarn,
		Info:       text,
	})
}
func (d *DeviceLogger) Warnf(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeWarn,
		Info:       fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Error(text string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeError,
		Info:       text,
	})
}
func (d *DeviceLogger) Errorf(text string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeError,
		Info:       fmt.Sprintf(text, val...),
	})
}

func (d *DeviceLogger) Send(info string, msg string) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeSend,
		Message:    msg,
		Info:       info,
	})
}
func (d *DeviceLogger) Sendf(info string, msg string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeSend,
		Info:       fmt.Sprintf(info, val...),
		Message:    msg,
	})
}

func (d *DeviceLogger) Receive(info string, msg string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeRecevice,
		Message:    msg,
		Info:       info,
	})
}
func (d *DeviceLogger) Receivef(info string, msg string, val ...interface{}) {
	repository.GetDeviceLogRepository().Add(&repository.DeviceLog{
		DeviceName: d.deviceName,
		Theme:      model.ThemeRecevice,
		Info:       fmt.Sprintf(info, val...),
		Message:    msg,
	})
}

func (d *DeviceLogger) Logs(theme model.Theme) ([]*model.DeviceLog, error) {
	logs, err := repository.GetDeviceLogRepository().Query(d.deviceName, theme)
	if err != nil {
		fmt.Println("query logs failed", err)
		return nil, err
	}
	out := make([]*model.DeviceLog, len(logs))
	for i := range logs {
		out[i] = &model.DeviceLog{
			DeviceName: logs[i].DeviceName,
			Theme:      logs[i].Theme,
			Info:       logs[i].Info,
			Message:    logs[i].Message,
			CreatedAt:  logs[i].CreatedAt,
		}
	}
	return out, nil
}

func (d *DeviceLogger) RangeLogs(theme model.Theme, start, end time.Time) ([]*model.DeviceLog, error) {
	logs, err := repository.GetDeviceLogRepository().QueryRange(d.deviceName, theme, start, end)
	if err != nil {
		fmt.Println("query logs failed", err)
		return nil, err
	}
	out := make([]*model.DeviceLog, len(logs))
	for i := range logs {
		out[i] = &model.DeviceLog{
			DeviceName: logs[i].DeviceName,
			Theme:      logs[i].Theme,
			Info:       logs[i].Info,
			Message:    logs[i].Message,
			CreatedAt:  logs[i].CreatedAt,
		}
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
