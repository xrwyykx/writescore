package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"writescore/app"
	"writescore/data/db/user"
	"writescore/models/co"
	"writescore/models/dto"
)

//	type UploadEssayByHandMap struct {
//		//UserID       int64  `json:"userId" gorm:"column:user_id;not null"`
//		Title        string `json:"title" gorm:"column:title;not null"`
//		Content      string `json:"content" gorm:"column:content"`
//		LanguageType int    `json:"languageType" gorm:"column:language_type;not null"`
//		UploadMethod int    `json:"uploadMethod" gorm:"column:upload_method"`
//	}
//
// 上传作文然后需要返回评分结果的这个时候需要调用到ai来评分的
func UploadEssayByHand(c *gin.Context) {
	userId := app.GetUserId(c)
	var data dto.UploadEssayByHandMap
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	if err := user.UploadEssayByHand(c, data, userId); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("作文上传失败"+err.Error()))

	}
}
