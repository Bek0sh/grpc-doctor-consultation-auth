package models

import "github.com/Bek0sh/online-market-auth/pkg/proto"

type User struct {
	Id          int
	Name        string
	Username    string
	PhoneNumber string
	Password    string
}

type RegisterUser struct {
	Username        string
	Surname         string
	PhoneNumber     string
	Password        string
	ConfirmPassword string
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
	PhoneNumber string
	Password    string
}

func (s *SignInUser) ToProto() *proto.SignInRequest {
	return &proto.SignInRequest{
		PhoneNumber: s.PhoneNumber,
		Password:    s.Password,
	}
}
