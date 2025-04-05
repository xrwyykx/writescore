package main

import (
	"github.com/gin-gonic/gin"
	"writescore/global"
	"writescore/router"
)

func main() {
	//加载配置文件
	global.LoadConfig()
	gin.SetMode(gin.ReleaseMode)
	//加载redis
	//还未实现
	//连接数据库
	global.InitDB()
	//启动路由
	router.InitRouterAndStartServer()
}
