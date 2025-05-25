package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouterAndStartServer() {
	router := gin.Default()

	// 使用自定义CORS中间件
	router.Use(CorsHandler())

	// 设置路由
	root := router.Group(viper.GetString("http.path"))
	setCommonRouters(root)
	setUserRouters(root)

	// 启动服务器
	router.Run(":8099")
}

func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, OPTIONS, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Set-Cookie, Content-Length, Content-Range")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
