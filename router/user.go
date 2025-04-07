package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/user/controller"
)

func setUserRouters(root *gin.RouterGroup) {
	root.POST("/api/user/get-user-info", controller.GetUserInfo)
	root.POST("/api/user/update-user-info", controller.UpdateUserInfo)
	root.POST("/api/user/update-password", controller.UpdatePassword)
}
