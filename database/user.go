package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"synthori/ediary/m/models"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserDatabase interface {
	GetUserByID(id int) (models.GetUser, error)
	GetUserByUsername(username string) (models.GetUser, error)
	GetUsers(limit int, offset int) ([]models.GetUser, error)
	CreateUser(user models.UserCreateForm) (*models.GetUser, error)
	UpdateUser(user *models.UpdateUser) (*models.User, error)
	DeleteUserByID(id int) (*models.User, error)
	DeleteUserByUsername(username string) (*models.User, error)
}

type SQLiteUserDatabase struct {
	db *sql.DB
}

func NewSQLiteUserDatabase(db *sql.DB) *SQLiteUserDatabase {
	return &SQLiteUserDatabase{db: db}
}

var _ UserDatabase = (*SQLiteUserDatabase)(nil)

func (repo *SQLiteUserDatabase) GetUserByID(id int) (models.GetUser, error) {
	var user models.GetUser
	err := repo.db.QueryRow("select id, username, coalesce(email, ''), role from users where id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.GetUser{}, errors.New("user not found")
	}
	if err != nil {
		return models.GetUser{}, err
	}

	return user, nil
}

func (repo *SQLiteUserDatabase) GetUserByUsername(username string) (models.GetUser, error) {
	var user models.GetUser
	err := repo.db.QueryRow("select id, username, coalesce(email, ''), role from users where username=?", username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.GetUser{}, errors.New("user not found")
	}
	if err != nil {
		return models.GetUser{}, err
	}

	return user, nil
}

func (repo *SQLiteUserDatabase) GetUsers(limit int, offset int) ([]models.GetUser, error) {
	users := []models.GetUser{}

	rows, err := repo.db.Query("select id, username, coalesce(email, ''), role from users limit ? offset ?", limit, offset)
	if err != nil {
		return []models.GetUser{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("Warning", "close file: ", fmt.Sprintf("%s", err))
		}
	}()

	for rows.Next() {
		var user models.GetUser
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
		)

		if err != nil {
			return []models.GetUser{}, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return []models.GetUser{}, err
	}

	return users, nil
}

func (repo *SQLiteUserDatabase) CreateUser(user models.UserCreateForm) (*models.GetUser, error) {

	hash_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	result, err := repo.db.Exec("insert into users (username, password_hash, role) values (?, ?, ?)", user.Username, hash_password, user.Role)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.GetUser{
		ID:       int(id),
		Username: user.Username,
		Email:    "",
		Role:     user.Role,
	}, nil
}

func (repo *SQLiteUserDatabase) UpdateUser(user *models.UpdateUser) (*models.User, error) {
	return nil, nil
}

func (repo *SQLiteUserDatabase) DeleteUserByID(id int) (*models.User, error) {
	return nil, nil
}

func (repo *SQLiteUserDatabase) DeleteUserByUsername(username string) (*models.User, error) {
	return nil, nil
}
