package db

import (
	"database/sql"
	"fmt"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var database *sql.DB

func Connect(cfg *config.Config) {
	conn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		cfg.Postgre.DbUsername,
		cfg.Postgre.DbPassword,
		cfg.Postgre.DbPort,
		cfg.Postgre.DbName,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		logrus.Fatal("failed to connect to db, error: ", err)
	}

	database = db
}

func GetDB() *sql.DB {
	return database
}
