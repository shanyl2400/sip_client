package repository

import (
	"encoding/json"
	"sipsimclient/errors"
	"sync"

	"github.com/boltdb/bolt"
)

type User struct {
	Name     string
	Password string
	Role     string
}

type UserRepository interface {
	Put(u *User) error
	Get(name string) (*User, error)
	Delete(name string) error

	All() ([]*User, error)
}

type BoltUserRepository struct {
}

func (b *BoltUserRepository) Put(u *User) error {
	err := Get().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UsersBucket))
		userJSON, err := json.Marshal(u)
		if err != nil {
			return err
		}
		err = b.Put([]byte(u.Name), userJSON)
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

func (b *BoltUserRepository) Get(name string) (*User, error) {
	user := new(User)
	err := Get().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UsersBucket))
		val := b.Get([]byte(name))
		if val == nil {
			return errors.ErrUserNotExists
		}
		err := json.Unmarshal(val, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (b *BoltUserRepository) Delete(name string) error {
	err := Get().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UsersBucket))
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
func (b *BoltUserRepository) All() ([]*User, error) {
	out := make([]*User, 0)
	err := Get().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UsersBucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			user := new(User)
			err := json.Unmarshal(v, user)
			if err != nil {
				return err
			}
			out = append(out, user)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

var (
	_userRepository     UserRepository
	_userRepositoryOnce sync.Once
)

func GetUserRepository() UserRepository {
	_userRepositoryOnce.Do(func() {
		_userRepository = new(BoltUserRepository)
	})
	return _userRepository
}
