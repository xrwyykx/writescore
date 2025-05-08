package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/user/controller"
)

func setUserRouters(root *gin.RouterGroup) {
	//个人信息相关
	root.POST("/api/user/get-user-info", controller.GetUserInfo)
	root.POST("/api/user/update-user-info", controller.UpdateUserInfo)
	root.POST("/api/user/update-password", controller.UpdatePassword)

	//文章相关
	root.POST("/api/user/rating-essay", controller.RantingEssay)         //给文章进行评分
	root.POST("/api/user/get-user-essay", controller.GetEssay)           //获取自己发布的文章
	root.POST("/api/user/get-essay-details", controller.GetEssayDetails) //获取文章详情，包括评分信息等等
}
