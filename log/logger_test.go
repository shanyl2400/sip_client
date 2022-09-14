package log

import (
	"fmt"
	"testing"

	"github.com/gookit/color"
)

func TestColor(t *testing.T) {
	color.Info.Tips("Info tips message")
	color.Notice.Tips("Notice tips message")
	color.Warn.Tips("Notice tips message")
	color.Error.Tips("Error tips message")
	color.Secondary.Tips("Secondary tips message")
	color.Danger.Tips("Secondary tips message")

	red := color.FgRed.Render
	green := color.FgGreen.Render
	// bgred := color.BgRed.Render

	warnRender := color.Error.Render
	fmt.Printf("%s line %s library %v\n", red("Command"), green("color"), warnRender("ERROR"))
}

func TestLog(t *testing.T) {
	SetConfig(LogConfig{
		DisableColor: false,
		LogLevel:     InfoLevel,
	})
	Debug("This is a debug log")
	Debugf("This is a debugf log, value: %v", 1)
	Info("This is a info log")
	Infof("This is a infof log, value: %v", 2)

	Warn("This is a warn log")
	Warnf("This is a warnf log, value: %v", 3)
	Error("This is a error log")
	Errorf("This is a errorf log, value: %v", 4)

	Fatal("This is a fatal log")
	Fatalf("This is a fatalf log, value: %v", 5)
}
