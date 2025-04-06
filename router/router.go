package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouterAndStartServer() {
	router := gin.Default()
	root := router.Group(viper.GetString("http.path"))
	setCommonRouters(root)
	setUserRouters(root)
	router.Run(":8099")
}
