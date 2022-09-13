package repository

import (
	"sipsimclient/config"
	"sipsimclient/errors"
	"testing"
)

func TestPutDevice(t *testing.T) {
	config.Set(&config.Config{
		BoltDBPath: "./blot.db",
	})
	Init()
	defer Close()
	repo := GetDeviceRepository()

	d1 := &Device{
		Name:     "d1",
		Password: "123456",
		Protocol: "tcp",
	}
	d2 := &Device{
		Name:     "d2",
		Password: "123456",
		Protocol: "udp",
	}

	err := repo.Put(d1)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Put(d2)
	if err != nil {
		t.Fatal(err)
	}

	d1cpy, err := repo.Get(d1.Name)
	if err != nil {
		t.Fatal(err)
	}
	if d1cpy.Name != d1.Name || d1cpy.Password != d1.Password || d1cpy.Protocol != d1.Protocol {
		t.Fatal("device d1 not equal")
	}

	devices, err := repo.All()
	if err != nil {
		t.Fatal(err)
	}
	for i := range devices {
		t.Logf("%#v\n", devices[i])
	}

	err = repo.Delete(d1.Name)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Get(d1.Name)
	if err != errors.ErrDeviceNotExists {
		t.Fatal("should be deleted")
	}

	devices, err = repo.All()
	if err != nil {
		t.Fatal(err)
	}
	for i := range devices {
		t.Logf("%#v\n", devices[i])
	}
}
