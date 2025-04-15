package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/comon/controller"
)

func setCommonRouters(root *gin.RouterGroup) {
	root.POST("/api/common/get-upload-token", controller.GetUploadToken)     //获取上传凭证
	root.POST("/api/common/restore-image-info", controller.RestoreImageInfo) //上传后返回的信息返回添加到数据库
	root.POST("/api/common/register", controller.Register)                   //注册
	root.POST("/api/common/login", controller.Login)                         //登录
	root.POST("/api/common/shibie", controller.RecognizeText)                //识别文字
	//root.POST("/api/common/get-access-token", controller.GetAccessToken)     //获取新的
}
