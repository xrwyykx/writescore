package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouterAndStartServer() {
	router := gin.Default()
	root := router.Group(viper.GetString("http.path"))
	setCommonRouters(root)
	setUserRouters(root)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true //允許所有域名
	router.Use(cors.New(corsConfig))
	router.Run(":8099")
}
