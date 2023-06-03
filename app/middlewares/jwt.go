package middlewares

import (
	"net/http"
	"strings"

	"github.com/fajaaro/dbo/app"
	"github.com/fajaaro/dbo/app/controllers"
	"github.com/fajaaro/dbo/app/models"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := models.JsonResponse{Success: true}

		if len(strings.Split(c.Request.Header.Get("Authorization"), " ")) != 2 {
			errorMsg := "invalid token"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		accessToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]

		user, err := controllers.ValidateAccessToken(accessToken, app.GetDb())
		if err != nil {
			errorMsg := err.Error()
			res.Success = false
			res.Error = &errorMsg

			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
