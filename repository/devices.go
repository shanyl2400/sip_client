package repository

import (
	"encoding/json"
	"sipsimclient/errors"
	"sync"

	"github.com/boltdb/bolt"
)

type Device struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
}

type DeviceRepository interface {
	Put(d *Device) error
	Get(name string) (*Device, error)
	All() ([]*Device, error)

	Delete(name string) error
}

type BoltDeviceRepository struct{}

func (b *BoltDeviceRepository) Put(d *Device) error {
	err := Get().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DevicesBucket))
		deviceJSON, err := json.Marshal(d)
		if err != nil {
			return err
		}
		err = b.Put([]byte(d.Name), deviceJSON)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDeviceRepository) Get(name string) (*Device, error) {
	device := new(Device)
	err := Get().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DevicesBucket))
		val := b.Get([]byte(name))
		if val == nil {
			return errors.ErrDeviceNotExists
		}
		err := json.Unmarshal(val, device)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (b *BoltDeviceRepository) All() ([]*Device, error) {
	out := make([]*Device, 0)
	err := Get().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DevicesBucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			device := new(Device)
			err := json.Unmarshal(v, device)
			if err != nil {
				return err
			}
			out = append(out, device)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (b *BoltDeviceRepository) Delete(name string) error {
	err := Get().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DevicesBucket))
		err := b.Delete([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

var (
	_deviceRepository     DeviceRepository
	_deviceRepositoryOnce sync.Once
)

func GetDeviceRepository() DeviceRepository {
	_deviceRepositoryOnce.Do(func() {
		_deviceRepository = new(BoltDeviceRepository)
	})
	return _deviceRepository
}
