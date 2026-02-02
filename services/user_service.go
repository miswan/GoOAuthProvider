package services

import (
	"errors"
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
)

type UserService struct {
	store storage.Storage
}

func NewUserService(store storage.Storage) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Register(req *models.UserRegister) error {
	if s.store.GetUserByUsername(req.Username) != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	return s.store.StoreUser(user)
}

func (s *UserService) Login(req *models.UserLogin) (*models.User, error) {
	user := s.store.GetUserByUsername(req.Username)
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GetUser(id uint) *models.User {
	return s.store.GetUser(id)
}