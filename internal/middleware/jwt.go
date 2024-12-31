package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/configs"
	"github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/app"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, app.NewErrorResponse("Authorization header is required", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(401, app.NewErrorResponse("Invalid token format", nil))
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Config.JWT_SECRET), nil
		})

		if err != nil {
			c.JSON(401, app.NewErrorResponse("Unauthorized", nil))
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, app.NewErrorResponse("Invalid token", nil))
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func VerifyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(403, app.NewErrorResponse("You do not have permission to access this resource", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
