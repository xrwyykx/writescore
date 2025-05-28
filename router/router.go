package router

//import "C"
import (
	"log"

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
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Set-Cookie, Content-Length, Content-Range")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			// 打印请求头，用于调试
			log.Printf("OPTIONS 请求头: %v", c.Request.Header)
			c.AbortWithStatus(204)
			return
		}

		// 打印请求信息，用于调试
		log.Printf("请求信息: Method=%s, Path=%s, Cookie=%v", c.Request.Method, c.Request.URL.Path, c.Request.Cookies())

		c.Next()
	}
}
