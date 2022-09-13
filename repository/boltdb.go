package repository

import (
	"log"
	"sipsimclient/config"
	"sync"

	"github.com/boltdb/bolt"
)

var (
	_db     *bolt.DB
	_dbOnce sync.Once
)

func Close() {
	_db.Close()
}

func Get() *bolt.DB {
	_dbOnce.Do(func() {
		var err error
		_db, err = bolt.Open(config.Get().BoltDBPath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
	})
	return _db
}

func Init() {
	Get().Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(DevicesBucket))
		tx.CreateBucketIfNotExists([]byte(DeviceLogBucket))
		tx.CreateBucketIfNotExists([]byte(UsersBucket))
		return nil
	})
}
