package repository

import (
	"sipsimclient/config"
	"sipsimclient/errors"
	"testing"
)

func TestPutUser(t *testing.T) {
	config.Set(&config.Config{
		BoltDBPath: "./blot.db",
	})
	Init()
	defer Close()
	repo := GetUserRepository()

	user1 := &User{
		Name:     "user1",
		Password: "123456",
	}
	user2 := &User{
		Name:     "user2",
		Password: "2222",
	}

	err := repo.Put(user1)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Put(user2)
	if err != nil {
		t.Fatal(err)
	}

	d1cpy, err := repo.Get(user1.Name)
	if err != nil {
		t.Fatal(err)
	}
	if d1cpy.Name != user1.Name || d1cpy.Password != user1.Password {
		t.Fatal("user not equal")
	}

	err = repo.Delete(user1.Name)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Get(user1.Name)
	if err != errors.ErrUserNotExists {
		t.Fatal(err)
	}
}
