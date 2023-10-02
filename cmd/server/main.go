package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/repository"
	"github.com/Bek0sh/online-market-auth/internal/service"
	"github.com/Bek0sh/online-market-auth/pkg/db"
	"github.com/Bek0sh/online-market-auth/pkg/proto"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var srv proto.AuthUserServer
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
	srv = service.NewService(repo, cfg, jwt)
}

func main() {

	go grpcRun()
	grpcGateWayRun()
	time.Sleep(time.Second)
}

func grpcRun() {
	grpcServer := grpc.NewServer()

	proto.RegisterAuthUserServer(grpcServer, srv)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		logrus.Fatal(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		logrus.Fatal(err)
	}
}

func grpcGateWayRun() {
	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := proto.RegisterAuthUserHandlerServer(ctx, grpcMux, srv)
	if err != nil {
		logrus.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatal(err)
	}

	if err = http.Serve(listener, mux); err != nil {
		logrus.Fatal(err)
	}
}
