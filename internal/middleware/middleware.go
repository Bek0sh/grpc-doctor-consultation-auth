package middleware

import (
	"net/http"
	"strings"

	"github.com/Bek0sh/online-market-auth/pkg/token"
	"github.com/gin-gonic/gin"
)

func CheckUser(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		token, err := ctx.Cookie("access_token")
		header := ctx.GetHeader("Autorization")

		fields := strings.Fields(header)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[0]
		} else if err == nil {
			accessToken = token
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  "fail",
					"message": "you are not logged in",
				},
			)
			return
		}

		payload, err := maker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  "fail",
					"message": err.Error(),
				},
			)
			return
		}

		ctx.Set("user_id", payload.Id)
		ctx.Set("user_role", payload.UserRole)

		ctx.Next()
	}
}

func CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("user_role")

		if role.(string) != "DOCTOR" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  "fail",
					"message": "Only Doctors have permission for this action",
				},
			)
			return
		}

		c.Next()
	}
}
