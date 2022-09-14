package devices

import (
	"sipsimclient/config"
	"sipsimclient/model"
	"sipsimclient/repository"
	"testing"
)

func TestLog(t *testing.T) {
	config.Set(&config.Config{
		BoltDBPath: "./blot.db",
	})
	repository.Init()
	defer repository.Close()
	logger, err := NewLogger("device1")
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("设备已启动")
	logger.Info("设备运行中")
	logger.Infof("设备运行中%v", 123465)

	t.Log(logger.Logs(model.ThemeAll))
}
