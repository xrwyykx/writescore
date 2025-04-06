package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/comon/controller"
)

func setCommonRouters(root *gin.RouterGroup) {
	root.POST("/api/common/get-upload-token", controller.GetUploadToken)
	root.POST("api/common/register", controller.Register) //注册
	root.POST("api/common/login", controller.Login)       //登录
}
