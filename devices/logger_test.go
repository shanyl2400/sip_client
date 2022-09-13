package devices

import "testing"

func TestLog(t *testing.T) {
	logger, err := NewLogger("device1")
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("设备已启动")
	logger.Info("设备运行中")
	logger.Infof("设备运行中%v", 123465)
}
