package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/Bek0sh/online-market-auth/pkg/proto"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/golang/protobuf/ptypes"
)

type GrpcService struct {
	repo Repository
	jwt  token.Maker
	proto.UnimplementedUserInfoServer
}

func NewGrpcService(jwt token.Maker, repo Repository) proto.UserInfoServer {
	return &GrpcService{jwt: jwt, repo: repo}
}

func (g *GrpcService) GetCurrentUser(ctx context.Context, req *proto.Empty) (*proto.GetUserResponse, error) {

	payload, err := g.jwt.VerifyToken(access_token)

	if err != nil {
		return nil, err
	}

	user, err := g.repo.GetUserById(payload.Id)

	if err != nil {
		return nil, err
	}

	resp := modelToProto(user, payload.Id)
	return resp, nil
}

func modelToProto(user *models.UserResponse, id int) *proto.GetUserResponse {
	protoTimestamp, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		log.Fatalf("Error converting time to Protobuf Timestamp: %v", err)
	}

	protobufTimeString := protoTimestamp.AsTime().Format(time.RFC3339)

	resp := &proto.GetUserResponse{
		Id:          int32(id),
		Name:        user.Username,
		Surname:     user.Surname,
		PhoneNumber: user.PhoneNumber,
		UserRole:    user.UserRole,
		Email:       user.Email,
		CreatedAt:   protobufTimeString,
	}
	return resp
}

func (g *GrpcService) CheckToken(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	_, err := g.jwt.VerifyToken(access_token)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (g *GrpcService) CheckRole(context.Context, *proto.Empty) (*proto.Empty, error) {
	payload, err := g.jwt.VerifyToken(access_token)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(payload.UserRole) != "doctor" {
		return nil, fmt.Errorf("this action only for doctors, access denied")
	}

	return &proto.Empty{}, nil
}

func (g *GrpcService) GetUserById(ctx context.Context, req *proto.GetUserByIdRequest) (*proto.GetUserResponse, error) {
	user, err := g.repo.GetUserById(int(req.GetId()))

	if err != nil {
		return nil, err
	}

	resp := modelToProto(user, int(req.GetId()))

	return resp, nil
}
