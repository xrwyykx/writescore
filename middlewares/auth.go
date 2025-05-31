package middlewares

import (
	"net/http"
	"writescore/models/co"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查JWT token
		tokenstring := c.GetHeader("Authorization")
		if tokenstring != "" {
			claims, err := PraseToken(tokenstring)
			if err == nil {
				c.Set("user_id", claims.UserId)
				c.Next()
				return
			}
		}

		// 2. 检查SESSION cookie
		sessionCookie, err := c.Cookie("SESSION")
		if err != nil {
			c.JSON(http.StatusUnauthorized, co.BadRequest("未登录或会话已过期"))
			c.Abort()
			return
		}

		// 3. 验证SESSION
		userId, err := validateSession(sessionCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, co.BadRequest("会话无效"))
			c.Abort()
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
