package repository_test

import (
	"os"
	"testing"

	"github.com/Bek0sh/online-market-auth/internal/config"
	"github.com/Bek0sh/online-market-auth/internal/repository"
	"github.com/Bek0sh/online-market-auth/internal/service"
	"github.com/Bek0sh/online-market-auth/pkg/db"
)

var repo service.Repository

func TestMain(m *testing.M) {
	cfg := config.GetConfig()
	db.Connect(cfg)
	repo = repository.NewRepository(db.GetDB())

	os.Exit(m.Run())
}
