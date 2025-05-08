package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"writescore/app"
	"writescore/data/db/user"
	"writescore/models"
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
// // 上传作文然后需要返回评分结果的这个时候需要调用到ai来评分的
//
//	func UploadEssayByHand(c *gin.Context) {
//		userId := app.GetUserId(c)
//		var data dto.UploadEssayByHandMap
//		if err := c.ShouldBindJSON(&data); err != nil {
//			c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
//			return
//		}
//		if err := user.UploadEssayByHand(c, data, userId); err != nil {
//			c.JSON(http.StatusBadRequest, co.BadRequest("作文上传失败"+err.Error()))
//
//		}
//	}
//
// 已经可以实现将照片转化为文字，接下来就是评分阶段
func RantingEssay(c *gin.Context) {
	userId := app.GetUserId(c)
	var param dto.RatingEssayMap
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}

	data, err := user.RantingEssay(c, param, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("评分失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, co.Success("评分成功", data))

}

func GetEssay(c *gin.Context) {
	//按找提交时间起止来寻找，还有总评分范围，文章标题
	userId := app.GetUserId(c)
	var param dto.GetEssayMap
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	data, total, err := user.GetEssay(c, userId, param)
	if err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("获取我发布的文章失败,"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, co.Success("获取我发布的文章成功", gin.H{
		"data":  data,
		"total": total,
	}))

}

func GetEssayDetails(c *gin.Context) {
	userId := app.GetUserId(c)
	var param models.IdMap
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	data, err := user.GetEssayDetails(c, userId, param.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("获取作文详情失败"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, co.Success("获取作文详情成功", data))

}
