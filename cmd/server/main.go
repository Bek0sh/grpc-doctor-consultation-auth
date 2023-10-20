package main

import (
	"log"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/handler"
	"github.com/Bek0sh/online-market-auth/internal/repository"
	"github.com/Bek0sh/online-market-auth/internal/router"
	"github.com/Bek0sh/online-market-auth/internal/service"
	"github.com/Bek0sh/online-market-auth/pkg/db"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var handlers handler.Handler
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
	handlers = *handler.NewHandler(srv)
}

func main() {
	r := gin.Default()

	router.AuthRouter(r, jwt, handlers)

	log.Fatal(r.Run(":8080"))
}
