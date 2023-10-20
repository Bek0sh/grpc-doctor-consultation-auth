package router

import (
	"github.com/Bek0sh/online-market-auth/internal/handler"
	"github.com/Bek0sh/online-market-auth/internal/middleware"
	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine, maker token.Maker, handlers handler.Handler) {
	r.POST("/v1/auth/register", handlers.RegisterUser())
	r.POST("/v1/auth/sign-in", handlers.SignIn())
	r.GET("/v1/auth/profile", middleware.CheckUser(maker), handlers.Profile())
}
