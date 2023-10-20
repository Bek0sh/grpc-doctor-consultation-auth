package repository

import (
	"database/sql"

	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/Bek0sh/online-market-auth/internal/service"
	"github.com/sirupsen/logrus"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) service.Repository {
	return &repo{db: db}
}

func (r *repo) CreateUser(userInput *models.RegisterUser) (int, error) {
	var id int

	query := "INSERT INTO users (name, surname, phone_number,user_role, created_at ,hashed_password) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	err := r.db.QueryRow(query, userInput.Username, userInput.Surname, userInput.PhoneNumber, userInput.UserRole, userInput.CreatedAt, userInput.Password).Scan(&id)
	if err != nil {
		logrus.Error("failed to insert user into database, error: ", err)
		return id, err
	}

	return id, nil
}

func (r *repo) GetUserById(id int) (*models.UserResponse, error) {
	response := models.UserResponse{}

	query := "SELECT name, surname, user_role, created_at, phone_number FROM users WHERE id=$1"

	err := r.db.QueryRow(query, &id).Scan(&response.Username, &response.Surname, &response.UserRole, &response.CreatedAt, &response.PhoneNumber)

	if err != nil {
		logrus.Errorf("failed to find user with id=%d, error: %v", id, err)
		return &models.UserResponse{}, err
	}

	return &response, nil
}

func (r *repo) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	response := models.User{}

	query := "SELECT id, name, surname, phone_number, user_role, created_at, hashed_password FROM users WHERE phone_number=$1"

	err := r.db.QueryRow(query, &phoneNumber).Scan(&response.Id, &response.Username, &response.Surname, &response.PhoneNumber, &response.UserRole, &response.CreatedAt, &response.Password)
	if err != nil {
		logrus.Errorf("failed to find user with phone number=%s, error: %v", phoneNumber, err)
		return &models.User{}, err
	}
	return &response, nil
}

func (r *repo) UpdateUser(userInput *models.User) (*models.UserResponse, error) {
	response := models.UserResponse{}

	query := "UPDATE users SET name=$1, surname=$2 WHERE id=$3 RETURNING id, name, surname, phone_number"
	err := r.db.QueryRow(query, userInput.Username, userInput.Surname, userInput.Id).Scan(&response.Id, &response.Username, &response.Surname, &response.PhoneNumber)

	if err != nil {
		logrus.Errorf("failed to update user with id=%d, error: %v", userInput.Id, err)
		return &models.UserResponse{}, err
	}

	return &response, err
}

func (r *repo) DeleteUser(id int) error {

	query := "DELETE FROM users WHERE id=$1"

	_, err := r.db.Exec(query, id)
	if err != nil {
		logrus.Errorf("failed to delete user with id=%d, error: %v", id, err)
		return err
	}
	return nil
}
