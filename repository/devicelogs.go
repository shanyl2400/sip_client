package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

const (
	ThemeSend     = "send"
	ThemeRecevice = "receive"

	ThemeAll         = "all"
	ThemeSendRecv    = "send/recv"
	ThemeTransaction = "transaction"
)

type Theme string

func (t Theme) filter(t2 Theme) bool {
	if t == t2 || t == ThemeAll {
		return true
	}
	if t == ThemeSendRecv && (t2 == ThemeSend || t2 == ThemeRecevice) {
		return true
	}

	return false
}

type DeviceLog struct {
	DeviceName string
	Theme      Theme
	Message    string
	CreatedAt  time.Time
}

type DeviceLogRepository interface {
	Add(d *DeviceLog) error
	Query(name string, theme Theme) ([]*DeviceLog, error)
	DeleteAll(name string) error
}

type BoltDeviceLogRepository struct{}

func (b *BoltDeviceLogRepository) Add(d *DeviceLog) error {
	d.CreatedAt = time.Now()

	err := Get().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceLogBucket))
		logJSON, err := json.Marshal(d)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%v:%v", d.DeviceName, time.Now())
		err = b.Put([]byte(key), logJSON)
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

func (b *BoltDeviceLogRepository) Query(name string, theme Theme) ([]*DeviceLog, error) {
	out := make([]*DeviceLog, 0)
	err := Get().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceLogBucket))

		c := b.Cursor()

		for k, v := c.Seek([]byte(name)); k != nil && bytes.HasPrefix(k, []byte(name)); k, v = c.Next() {
			log := new(DeviceLog)
			err := json.Unmarshal(v, log)
			if err != nil {
				return err
			}
			if theme.filter(log.Theme) {
				out = append(out, log)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (b *BoltDeviceLogRepository) DeleteAll(name string) error {
	err := Get().Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DeviceLogBucket))

		c := b.Cursor()

		for k, _ := c.Seek([]byte(name)); k != nil && bytes.HasPrefix(k, []byte(name)); k, _ = c.Next() {
			b.Delete(k)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

var (
	_deviceLogRepository     DeviceLogRepository
	_deviceLogRepositoryOnce sync.Once
)

func GetDeviceLogRepository() DeviceLogRepository {
	_deviceLogRepositoryOnce.Do(func() {
		_deviceLogRepository = new(BoltDeviceLogRepository)
	})
	return _deviceLogRepository
}
