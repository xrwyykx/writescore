package router

import (
	"github.com/gin-gonic/gin"
	"writescore/app/comon/controller"
)

func setCommonRouters(root *gin.RouterGroup) {
	root.GET("/api/common/get-upload-token", controller.GetUploadToken)
}
