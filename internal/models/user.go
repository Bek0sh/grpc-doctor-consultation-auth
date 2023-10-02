package models

import "github.com/Bek0sh/online-market-auth/pkg/proto"

type User struct {
	Id          int
	Username    string
	Surname     string
	PhoneNumber string
	UserType    string
	Password    string
}

type UserResponse struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Surname     string `json:"surname"`
	UserType    string `json:"user_type"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterUser struct {
	Username        string `json:"username"`
	Surname         string `json:"surname"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_passsword"`
}

func (r *RegisterUser) ToProto() *proto.RegisterUserRequest {
	return &proto.RegisterUserRequest{
		Name:            r.Username,
		Surname:         r.Surname,
		PhoneNumber:     r.PhoneNumber,
		Password:        r.Password,
		ConfirmPassword: r.ConfirmPassword,
	}
}

type SignInUser struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (s *SignInUser) ToProto() *proto.SignInRequest {
	return &proto.SignInRequest{
		PhoneNumber: s.PhoneNumber,
		Password:    s.Password,
	}
}
