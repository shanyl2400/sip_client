package repository

import (
	"fmt"
	"sipsimclient/config"
	"sipsimclient/model"
	"testing"
	"time"
)

func TestPutLogs(t *testing.T) {
	config.Set(&config.Config{
		BoltDBPath: "./blot.db",
	})
	Init()
	defer Close()
	repo := GetDeviceLogRepository()

	err := repo.Add(&DeviceLog{
		DeviceName: "d1",
		Theme:      model.ThemeSend,
		Message:    "123123",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Add(&DeviceLog{
		DeviceName: "d1",
		Theme:      model.ThemeRecevice,
		Message:    "22222",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Add(&DeviceLog{
		DeviceName: "d2",
		Theme:      model.ThemeSend,
		Message:    "3333",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Add(&DeviceLog{
		DeviceName: "d1",
		Theme:      model.ThemeTransaction,
		Message:    "444444",
	})
	if err != nil {
		t.Fatal(err)
	}

	out, err := repo.Query("d1", model.ThemeAll)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-----0-----")
	for i := range out {
		t.Logf("%#v", out[i])
	}

	out, err = repo.Query("d1", model.ThemeSendRecv)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-----1-----")
	for i := range out {
		t.Logf("%#v", out[i])
	}

	out, err = repo.Query("d2", model.ThemeSend)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-----2-----")
	for i := range out {
		t.Logf("%#v", out[i])
	}

	out, err = repo.Query("d2", model.ThemeRecevice)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-----3-----")
	for i := range out {
		t.Logf("%#v", out[i])
	}

	err = repo.DeleteAll("d1")
	if err != nil {
		t.Fatal(err)
	}

	out, err = repo.Query("d1", model.ThemeAll)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-----4-----")
	for i := range out {
		t.Logf("%#v", out[i])
	}

	out, err = repo.Query("d2", model.ThemeAll)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("-----5-----")
	for i := range out {
		t.Logf("%#v", out[i])
	}
}

func TestQueryLogs(t *testing.T) {
	config.Set(&config.Config{
		BoltDBPath: "./sip_client.db",
	})
	Init()
	defer Close()
	repo := GetDeviceLogRepository()

	// _, err := repo.Query("d1", model.ThemeAll)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println("--------------")
	end := time.Now()
	start := end.Add(-time.Second * 60 * 10)
	fmt.Println("start:", start, "end:", end)
	out, err := repo.QueryRange("d1", model.ThemeAll, start, end)
	if err != nil {
		t.Fatal(err)
	}
	for i := range out {
		t.Log(out[i])
	}
}
