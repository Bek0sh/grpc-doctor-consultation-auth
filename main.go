package main

import (
	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetConfig()
	logrus.Print(cfg.Postgre.DbPassword)
}
