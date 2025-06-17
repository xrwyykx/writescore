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

// StreamRatingEssay 流式评分控制器
func StreamRatingEssay(c *gin.Context) {
	var param dto.RatingEssayMap
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}

	userId := app.GetUserId(c)
	if userId <= 0 {
		c.JSON(http.StatusBadRequest, co.BadRequest("未登录"))
		return
	}

	err := user.StreamRatingEssay(c, param, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("获取作文详情失败"+err.Error()))
		return
	}
}
