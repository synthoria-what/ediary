package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"synthori/ediary/m/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserDatabase interface {
	GetUserByID(id int) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUsers(limit int, offset int) ([]models.GetUser, error)
	CreateUser(user models.UserCreateForm) (*models.User, error)
	UpdateUser(id int, req models.UpdateUser) (*models.User, error)
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

func (repo *SQLiteUserDatabase) GetUserByID(id int) (models.User, error) {
	var user models.User
	// var CreatedAtStr string
	err := repo.db.QueryRow(`
		SELECT id, username, coalesce(email, ''), role, password_hash, created_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, errors.New("user not found")
	}
	if err != nil {
		return models.User{}, err
	}

	// user.CreatedAt, err = time.Parse(time.UnixDate, CreatedAtStr)
	return user, nil
}

func (repo *SQLiteUserDatabase) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := repo.db.QueryRow("select id, username, coalesce(email, ''), role, password_hash, created_at from users where username=?", username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, errors.New("user not found")
	}
	if err != nil {
		return models.User{}, err
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

func (repo *SQLiteUserDatabase) CreateUser(req models.UserCreateForm) (*models.User, error) {
	now := time.Now().UTC()
	hash_password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := repo.db.Exec("insert into users (username, password_hash, role, email, created_at) values (?, ?, ?, ?, ?)",
		req.Username, hash_password, req.Role, req.Email, now)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        int(id),
		Username:  req.Username,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: now,
	}, nil
}

func (repo *SQLiteUserDatabase) UpdateUser(id int, req models.UpdateUser) (*models.User, error) {
	fmt.Println("updating user in database")

	user, err := repo.GetUserByID(id)
	if err != nil {
		return &models.User{}, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	_, err = repo.db.Exec(`
		UPDATE users
		SET username = ?, role = ?, email = ?
		WHERE id = ?
	`, user.Username, user.Role, user.Email, id)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

func (repo *SQLiteUserDatabase) DeleteUserByID(id int) (*models.User, error) {
	user, err := repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	_, err = repo.db.Exec("delete from users where id = ?", user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *SQLiteUserDatabase) DeleteUserByUsername(username string) (*models.User, error) {
	user, err := repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	_, err = repo.db.Exec("delete from users where id = ?", user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
