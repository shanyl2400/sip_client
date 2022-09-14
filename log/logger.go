package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gookit/color"
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var (
	_logger *log.Logger
	_config LogConfig

	_levelRender = []func(a ...interface{}) string{
		color.Secondary.Render,
		color.Info.Render,
		color.Warn.Render,
		color.Error.Render,
		color.Danger.Render,
	}

	_levelName = []string{
		"DEBUG", "INFO", "WARN", "ERROR", "FATAL",
	}
)

type LogConfig struct {
	DisableColor bool
	LogLevel     int
}

func init() {
	_logger = log.New(os.Stdout, "", 0)
	_config = LogConfig{DisableColor: false, LogLevel: 0}
}

func SetConfig(c LogConfig) {
	_config = c
}

func Debug(text string) {
	Debugf(text)
}
func Debugf(text string, val ...interface{}) {
	text = sprintf(DebugLevel, text, val...)
	if text != "" {
		_logger.Print(text)
	}
}

func Info(text string) {
	Infof(text)
}
func Infof(text string, val ...interface{}) {
	text = sprintf(InfoLevel, text, val...)
	if text != "" {
		_logger.Print(text)
	}
}

func Warn(text string) {
	Warnf(text)
}
func Warnf(text string, val ...interface{}) {
	text = sprintf(WarnLevel, text, val...)
	if text != "" {
		_logger.Print(text)
	}
}

func Error(text string) {
	Errorf(text)
}
func Errorf(text string, val ...interface{}) {
	text = sprintf(ErrorLevel, text, val...)
	if text != "" {
		_logger.Print(text)
	}
}

func Fatal(text string) {
	Fatalf(text)
}
func Fatalf(text string, val ...interface{}) {
	text = sprintf(FatalLevel, text, val...)
	if text != "" {
		_logger.Print(text)
	}
}

func sprintf(level int, text string, val ...interface{}) string {
	now := time.Now().Format("2006/01/02 15:04:05.999")
	if _config.LogLevel > level {
		return ""
	}
	if _config.DisableColor {
		text = "[" + _levelName[level] + "] " + now + " " + text
	} else {
		text = _levelRender[level]("["+_levelName[level]+"]") + " " + now + " " + text
	}
	return fmt.Sprintf(text, val...)
}
