package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Bek0sh/online-market-auth/internal/models"
	"github.com/gin-gonic/gin"
)

type Service interface {
	SignInUser(ctx context.Context, req *models.SignInUser) (string, error)
	RegisterUser(ctx context.Context, req *models.RegisterUser) (int, error)
	GetProfile(ctx context.Context, id int) (*models.UserResponse, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInput models.RegisterUser

		if err := c.ShouldBindJSON(&userInput); err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"message": err.Error(),
					"status":  "fail",
				},
			)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id, err := h.service.RegisterUser(ctx, &userInput)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"message": err.Error(),
					"status":  "fail",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"id":     id,
				"status": "success",
			},
		)
	}
}

func (h *Handler) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInput models.SignInUser
		if err := c.ShouldBindJSON(&userInput); err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"message": err.Error(),
					"status":  "fail",
				},
			)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		token, err := h.service.SignInUser(ctx, &userInput)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"message": err.Error(),
					"status":  "fail",
				},
			)
			return
		}

		c.SetCookie("access_token", token, 20*60, "/", "localhost", false, true)

		c.JSON(
			http.StatusOK,
			gin.H{
				"your_token": token,
				"status":     "success",
			},
		)
	}
}

var userId int

func (h *Handler) Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		id := c.MustGet("user_id")

		userId = id.(int)
		user, err := h.service.GetProfile(ctx, userId)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"message": err.Error(),
					"status":  "fail",
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"current_user": user,
				"status":       "success",
			},
		)
	}
}
