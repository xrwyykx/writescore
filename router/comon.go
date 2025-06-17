package router

import (
	"writescore/app/comon/controller"

	"github.com/gin-gonic/gin"
)

func setCommonRouters(root *gin.RouterGroup) {
	root.POST("/api/common/get-upload-token", controller.GetUploadToken)                //获取上传凭证
	root.POST("/api/common/restore-image-info", controller.RestoreImageInfo)            //上传后返回的信息返回添加到数据库
	root.POST("/api/common/restore-multi-image-info", controller.RestoreMultiImageInfo) //上传多页图片并识别文字
	root.POST("/api/common/save-essay", controller.SaveEssay)                           //保存文章
	root.POST("/api/common/update-essay-content", controller.UpdateEssayContent)        //修改文章内容
	root.POST("/api/common/register", controller.Register)                              //注册
	root.POST("/api/common/login", controller.Login)                                    //登录
	root.POST("/api/common/logout", controller.Logout)                                  //退出登录
	root.POST("/api/common/shibie", controller.RecognizeText)                           //识别文字
	//root.POST("/api/common/get-access-token", controller.GetAccessToken)     //获取新的
}
