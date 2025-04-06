package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/user/controller"
)

func setUserRouters(root *gin.RouterGroup) {
	root.POST("/api/user/get-user-info", controller.GetUserInfo)
}
