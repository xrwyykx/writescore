package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func InitRouterAndStartServer() {
	router := gin.Default()
	root := router.Group(viper.GetString("http.path"))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true //允許所有域名
	log.Printf("CORS configuration: %+v\n", corsConfig)
	router.Use(cors.New(corsConfig))
	setCommonRouters(root)
	setUserRouters(root)
	router.Run(":8099")
}
