package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"writescore/models/co"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring := c.GetHeader("Authorization")
		if tokenstring == "" {
			c.JSON(http.StatusBadRequest, co.BadRequest("未获取到token"))
			c.Abort()
			return
		}
		claims, err := PraseToken(tokenstring)
		if err != nil {
			c.JSON(http.StatusBadRequest, co.BadRequest("token无效"))
			c.Abort()
			return
		}
		//若token有效，就将user_id藏在上下文中
		c.Set("user_id", claims.UserId)
		log.Printf("Token valid, user_id: %d\n", claims.UserId)
		c.Next()
	}
}
