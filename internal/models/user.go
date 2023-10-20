package models

import "time"

type User struct {
	Id          int
	Username    string
	Surname     string
	PhoneNumber string
	UserRole    string
	Password    string
	CreatedAt   time.Time
}

type UserResponse struct {
	Id          int    `json:"-"`
	Username    string `json:"username"`
	Surname     string `json:"surname"`
	UserRole    string `json:"user_role"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   time.Time
}

type RegisterUser struct {
	Username        string    `json:"username"`
	Surname         string    `json:"surname"`
	PhoneNumber     string    `json:"phone_number"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"created_at"`
	UserRole        string    `json:"user_role"`
	ConfirmPassword string    `json:"confirm_password"`
}

type SignInUser struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
