package devices

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	infoLevel  = "INFO"
	warnLeve   = "WARN"
	errLevel   = "ERR"
	fatalLevel = "FATAL"
	recvLevel  = "RECV"
	sendLevel  = "SEND"
)

type DeviceLogger struct {
	logger *log.Logger
	file   *os.File
}

func (d *DeviceLogger) Info(text string) {
	d.logger.Print(d.log(infoLevel, text))
}
func (d *DeviceLogger) Infof(text string, val ...interface{}) {
	d.logger.Printf(d.log(infoLevel, text), val...)
}

func (d *DeviceLogger) Warn(text string) {
	d.logger.Print(d.log(warnLeve, text))
}
func (d *DeviceLogger) Warnf(text string, val ...interface{}) {
	d.logger.Printf(d.log(warnLeve, text), val...)
}

func (d *DeviceLogger) Error(text string) {
	d.logger.Print(d.log(errLevel, text))
}
func (d *DeviceLogger) Errorf(text string, val ...interface{}) {
	d.logger.Printf(d.log(errLevel, text), val...)
}

func (d *DeviceLogger) Send(text string) {
	d.logger.Print(d.log(sendLevel, "\n"+text))
}
func (d *DeviceLogger) Sendf(text string, val ...interface{}) {
	d.logger.Printf(d.log(sendLevel, "\n"+text), val...)
}

func (d *DeviceLogger) Receive(text string, val ...interface{}) {
	d.logger.Printf(d.log(recvLevel, "\n"+text), val...)
}
func (d *DeviceLogger) Receivef(text string, val ...interface{}) {
	d.logger.Printf(d.log(recvLevel, "\n"+text), val...)
}

func (d *DeviceLogger) Fatal(text string) {
	d.logger.Fatal(d.log(fatalLevel, text))
}
func (d *DeviceLogger) Fatalf(text string, val ...interface{}) {
	d.logger.Printf(d.log(fatalLevel, text), val...)
}

func (d *DeviceLogger) log(level string, text string) string {
	return fmt.Sprintf("[%v] %v %v", level, time.Now().Format("2006/01/02 03:04:05.9999"), text)
}

func NewLogger(deviceName string) (*DeviceLogger, error) {
	logFileName := fmt.Sprintf("%v_%v.log", deviceName, time.Now().Format("20060102"))
	f, err := os.Create(logFileName)
	if err != nil {
		return nil, err
	}
	logger := log.New(f, "", 0)
	return &DeviceLogger{
		file:   f,
		logger: logger,
	}, nil
}
