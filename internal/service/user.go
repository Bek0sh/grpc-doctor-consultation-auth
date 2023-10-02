package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/Bek0sh/online-market-auth/pkg/proto"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/Bek0sh/online-market-auth/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	CreateUser(*models.RegisterUser) (int, error)
	GetUserById(int) (*models.UserResponse, error)
	GetUserByPhoneNumber(string) (*models.User, error)
	UpdateUser(*models.User) (*models.UserResponse, error)
	DeleteUser(int) error
}

type service struct {
	cfg  *config.Config
	repo Repository
	jwt  token.Maker
	proto.UnimplementedAuthUserServer
}

func NewService(repo Repository, cfg *config.Config, jwt token.Maker) proto.AuthUserServer {
	return &service{repo: repo, cfg: cfg, jwt: jwt}
}

func (s *service) RegisterUser(ctx context.Context, userInput *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	hashedPassword, err := utils.HashPassword(userInput.GetPassword())
	if err != nil {
		logrus.Error("Failed to hash password, error: ", err)
		return nil, err
	}
	createUser := &models.RegisterUser{
		Username:    userInput.GetName(),
		Surname:     userInput.GetSurname(),
		PhoneNumber: userInput.GetPhoneNumber(),
		Password:    hashedPassword,
	}

	id, err := s.repo.CreateUser(createUser)
	if err != nil {
		logrus.Error("failed to create user in service layer, error: ", err)
		return nil, err
	}

	response := &proto.RegisterUserResponse{
		Id: int32(id),
	}

	return response, err
}

func (s *service) SignInUser(ctx context.Context, userInput *proto.SignInRequest) (*proto.SignInResponse, error) {
	user, err := s.repo.GetUserByPhoneNumber(userInput.GetPhoneNumber())
	if err != nil {
		logrus.Errorf("failed to find user with phone-number=%s, error: %v", userInput.GetPhoneNumber(), err)
		return nil, err
	}

	if err = utils.ComparePassword(userInput.GetPassword(), user.Password); err != nil {
		logrus.Error("password is not correct, check it please, error: ", err)
		return nil, err
	}

	dur, _ := strconv.Atoi(s.cfg.JWT.AccessTokenDuration)

	access_token, err := s.jwt.CreateToken(user.Id, user.Username, time.Duration(time.Duration(dur).Minutes()))
	if err != nil {
		logrus.Error("failed to create access_token, error: ", err)
		return nil, err
	}

	return &proto.SignInResponse{AccessToken: access_token}, nil
}
