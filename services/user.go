package services

import (
	"synthori/ediary/m/database"
	"synthori/ediary/m/models"

	"golang.org/x/crypto/bcrypt"
)

type User interface {
	GetUserByID(id int) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	AuthenticateUser(username string, password string) (models.User, error)
	LoginUser(username string, password string) (string, error)
	GetUsers(limit int, offset int) ([]models.GetUser, error)
	CreateUser(user models.UserCreateForm) (*models.User, error)
	UpdateUser(user models.UpdateUser) (*models.User, error)
	DeleteUserByID(id int) (*models.User, error)
	DeleteUserByUsername(username string) (*models.User, error)
}

var _ User = (*UserService)(nil)

type UserService struct {
	repo database.UserDatabase
}

func NewUserService(repo database.UserDatabase) *UserService {
	return &UserService{repo: repo}
}

func (u UserService) GetUserByID(id int) (models.User, error) {
	user, err := u.repo.GetUserByID(id)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u UserService) GetUserByUsername(username string) (models.User, error) {
	user, err := u.repo.GetUserByUsername(username)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u UserService) AuthenticateUser(username string, password string) (models.User, error) {
	user, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u UserService) LoginUser(username string, password string) (string, error) {
	user, err := u.AuthenticateUser(username, password)
	if err != nil {
		return "", err
	}

	token, err := GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (u UserService) GetUsers(limit int, offset int) ([]models.GetUser, error) {
	users, err := u.repo.GetUsers(limit, offset)
	if err != nil {
		return []models.GetUser{}, err
	}

	return users, nil
}

func (u UserService) CreateUser(user models.UserCreateForm) (*models.User, error) {
	newUser, err := u.repo.CreateUser(user)
	if err != nil {
		return &models.User{}, nil
	}

	return newUser, nil
}

func (u UserService) UpdateUser(req models.UpdateUser) (*models.User, error) {
	return nil, nil
}

func (u UserService) DeleteUserByID(id int) (*models.User, error) {
	return nil, nil
}

func (u UserService) DeleteUserByUsername(username string) (*models.User, error) {
	return nil, nil
}
