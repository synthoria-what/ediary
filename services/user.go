package services

import (
	"synthori/ediary/m/database"
	"synthori/ediary/m/models"
)

type User interface {
	GetUserByID(id int) (models.GetUser, error)
	GetUserByUsername(username string) (models.GetUser, error)
	GetUsers(limit int, offset int) ([]models.GetUser, error)
	CreateUser(user models.User) (*models.GetUser, error)
	UpdateUser(user *models.UpdateUser) (*models.User, error)
	DeleteUserByID(id int) (*models.User, error)
	DeleteUserByUsername(username string) (*models.User, error)
}

type UserService struct {
	repo database.UserDatabase
}

func NewUserService(repo database.UserDatabase) *UserService {
	return &UserService{repo: repo}
}

func (u UserService) GetUserByID(id int) (models.GetUser, error) {
	user, err := u.repo.GetUserByID(id)

	if err != nil {
		return models.GetUser{}, err
	}

	return user, nil
}

func (u UserService) GetUsers(limit int, offset int) ([]models.GetUser, error) {
	users, err := u.repo.GetUsers(limit, offset)
	if err != nil {
		return []models.GetUser{}, err
	}

	return users, nil
}

func (u UserService) CreateUser(user models.UserCreateForm) (*models.GetUser, error) {
	newUser, err := u.repo.CreateUser(user)
	if err != nil {
		return &models.GetUser{}, nil
	}

	return newUser, nil
}
