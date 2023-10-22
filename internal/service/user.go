package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/handler"
	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/Bek0sh/online-market-auth/pkg/utils"
	"github.com/sirupsen/logrus"
)

// var currentUser models.UserResponse

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
}

func NewService(repo Repository, cfg *config.Config, jwt token.Maker) handler.Service {
	return &service{repo: repo, cfg: cfg, jwt: jwt}
}

func (s *service) RegisterUser(ctx context.Context, req *models.RegisterUser) (int, error) {

	violations := validateRegisterUser(req)
	if violations != nil {
		return 0, InvalidArgumentError(violations)
	}

	if req.Password != req.ConfirmPassword {
		logrus.Error("passwords do not match")
		return 0, fmt.Errorf("both of passwords must be equal and same, check them")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logrus.Error("Failed to hash password, error: ", err)
		return 0, err
	}

	createUser := &models.RegisterUser{
		Username:    req.Username,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   req.CreatedAt,
		UserRole:    req.UserRole,
		Email:       req.Email,
		Password:    hashedPassword,
	}

	id, err := s.repo.CreateUser(createUser)
	if err != nil {
		logrus.Error("failed to create user in service layer, error: ", err)
		return 0, err
	}

	return id, err
}

var access_token string

func (s *service) SignInUser(ctx context.Context, req *models.SignInUser) (string, error) {

	violations := validateSignInUser(req)
	if violations != nil {
		return "", InvalidArgumentError(violations)
	}

	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		logrus.Errorf("failed to find user with phone-number=%s, error: %v", req.PhoneNumber, err)
		return "", err
	}

	if err = utils.ComparePassword(req.Password, user.Password); err != nil {
		logrus.Error("password is not correct, check it please, error: ", err)
		return "", err
	}

	dur, _ := strconv.Atoi(s.cfg.JWT.AccessTokenDuration)

	access_token, err = s.jwt.CreateToken(user.Id, user.Username, user.UserRole, time.Duration(dur*int(time.Minute)))

	if err != nil {
		logrus.Error("failed to create access_token, error: ", err)
		return "", err
	}

	return access_token, nil
}

func (s *service) GetProfile(ctx context.Context, id int) (*models.UserResponse, error) {
	return s.repo.GetUserById(id)
}

func (s *service) GetToken() string {
	return access_token
}
