package user

import (
	"sipsimclient/errors"
	"sipsimclient/log"
	"sipsimclient/model"
	"sipsimclient/repository"
	"sipsimclient/tools"
)

const (
	resetPassword = "123456"
	adminName     = "admin"
	adminPassword = "123456"
)

type User interface {
	Init()
	Get(name string) (*model.User, error)
	Register(name, password string) error
	Login(name, password string) (string, error)
	List() ([]*model.User, error)

	Delete(name string) error
	UpdatePassword(name string, password string, newPassword string) error
	ResetPassword(name string) error
}

type user struct {
}

func (u *user) Init() {
	user, _ := repository.GetUserRepository().Get(adminName)
	if user != nil {
		return
	}
	err := repository.GetUserRepository().Put(&repository.User{
		Name:     adminName,
		Role:     model.RoleAdmin,
		Password: tools.Hash(adminPassword),
	})
	if err != nil {
		log.Warnf("Init admin failed, err: %v", err)
	}
}

func (u *user) Get(name string) (*model.User, error) {
	user, err := repository.GetUserRepository().Get(name)
	if err != nil {
		log.Infof("Get user failed, err: %v", err)
		return nil, err
	}
	return &model.User{
		Name:     user.Name,
		Password: user.Password,
		Role:     user.Role,
	}, nil
}
func (u *user) Register(name, password string) error {
	//check duplicate
	user, _ := repository.GetUserRepository().Get(name)
	if user != nil {
		log.Infof("User exists, name: %v", name)
		return errors.ErrUserExists
	}
	//put user
	return repository.GetUserRepository().Put(&repository.User{
		Name:     name,
		Password: tools.Hash(password),
		Role:     model.RoleCustomer,
	})
}
func (u *user) Login(name, password string) (string, error) {
	user, err := repository.GetUserRepository().Get(name)
	if err != nil {
		log.Infof("Get user failed, err: %v", err)
		return "", err
	}
	if tools.Hash(password) != user.Password {
		return "", errors.ErrIncorrectPassword
	}

	return tools.CreateToken(name, user.Role)
}

func (u *user) Delete(name string) error {
	user, _ := repository.GetUserRepository().Get(name)
	if user == nil {
		return nil
	}
	if user.Role == model.RoleAdmin {
		return errors.ErrRemoveAdmin
	}
	return repository.GetUserRepository().Delete(name)
}
func (u *user) List() ([]*model.User, error) {
	users, err := repository.GetUserRepository().All()
	if err != nil {
		log.Warnf("List users failed, err: %v", err)
		return nil, err
	}
	out := make([]*model.User, len(users))
	for i := range users {
		out[i] = &model.User{
			Name: users[i].Name,
			Role: users[i].Role,
		}
	}
	return out, nil
}
func (u *user) UpdatePassword(name string, password string, newPassword string) error {
	//check duplicate
	user, _ := repository.GetUserRepository().Get(name)
	if user == nil {
		log.Infof("user not exists, name: %v", name)
		return errors.ErrUserNotExists
	}
	if tools.Hash(password) != user.Password {
		return errors.ErrIncorrectPassword
	}
	//put user
	return repository.GetUserRepository().Put(&repository.User{
		Name:     user.Name,
		Password: tools.Hash(newPassword),
		Role:     user.Role,
	})
}
func (u *user) ResetPassword(name string) error {
	//check duplicate
	user, _ := repository.GetUserRepository().Get(name)
	if user == nil {
		log.Infof("user not exists, name: %v", name)
		return errors.ErrUserNotExists
	}
	//put user
	return repository.GetUserRepository().Put(&repository.User{
		Name:     user.Name,
		Password: tools.Hash(resetPassword),
		Role:     user.Role,
	})
}

func NewUser() User {
	return new(user)
}
