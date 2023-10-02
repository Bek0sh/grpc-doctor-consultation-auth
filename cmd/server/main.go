package main

import (
	"net"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/repository"
	"github.com/Bek0sh/online-market-auth/internal/service"
	"github.com/Bek0sh/online-market-auth/pkg/db"
	"github.com/Bek0sh/online-market-auth/pkg/proto"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var Service proto.AuthUserServer
var cfg *config.Config

func init() {
	cfg = config.GetConfig()
	db.Connect(cfg)
	database := db.GetDB()
	jwt, err := token.NewJwtMaker(cfg.JWT.SecretPassword)
	if err != nil {
		logrus.Error("failed to call NewJWTMaker, error: ", err)
	}
	repo := repository.NewRepository(database)
	Service = service.NewService(repo, cfg, jwt)
}

func main() {

	grpcServer := grpc.NewServer()
	proto.RegisterAuthUserServer(grpcServer, Service)
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatal(err)
	}

	if err = grpcServer.Serve(listen); err != nil {
		logrus.Fatal("failed to serve")
	}
}
