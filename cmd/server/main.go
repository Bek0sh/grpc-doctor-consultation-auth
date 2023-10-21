package main

import (
	"log"
	"net"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/handler"
	"github.com/Bek0sh/online-market-auth/internal/repository"
	"github.com/Bek0sh/online-market-auth/internal/router"
	"github.com/Bek0sh/online-market-auth/internal/service"
	"github.com/Bek0sh/online-market-auth/pkg/db"
	"github.com/Bek0sh/online-market-auth/pkg/proto"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var handlers handler.Handler
var rpc proto.UserInfoServer
var srv handler.Service
var cfg *config.Config
var jwt token.Maker

func init() {
	var err error
	cfg = config.GetConfig()
	db.Connect(cfg)
	database := db.GetDB()
	jwt, err = token.NewJwtMaker(cfg.JWT.SecretPassword)
	if err != nil {
		logrus.Error("failed to call NewJWTMaker, error: ", err)
	}
	repo := repository.NewRepository(database)
	srv = service.NewService(repo, cfg, jwt)
	rpc = service.NewGrpcService(jwt, repo)
	handlers = *handler.NewHandler(srv)
}

func main() {
	r := gin.Default()

	router.AuthRouter(r, jwt, handlers)

	go grpcServer()

	log.Fatal(r.Run(":8080"))
}

func grpcServer() {
	server := grpc.NewServer()

	reflection.Register(server)
	proto.RegisterUserInfoServer(server, rpc)

	liss, err := net.Listen("tcp", ":50051")
	if err != nil {
		logrus.Fatal(err)
	}

	if err = server.Serve(liss); err != nil {
		logrus.Fatal(err)
	}

}
