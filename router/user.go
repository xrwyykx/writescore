package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/user/controller"
)

func setUserRouters(root *gin.RouterGroup) {
	root.POST("/api/user/get-user-info", controller.GetUserInfo)
	root.POST("/api/user/update-user-info", controller.UpdateUserInfo)
	root.POST("/api/user/update-password", controller.UpdatePassword)
	root.POST("/api/user/upload-essay-by-hand", controller.UploadEssayByHand) //手写上传作文并返回评分结果
}
