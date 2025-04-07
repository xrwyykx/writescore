package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"writescore/app"
	"writescore/data/db/user"
	"writescore/models/co"
	"writescore/models/dto"
)

func GetUserInfo(c *gin.Context) {
	userId := app.GetUserId(c)
	data, err := user.GetUserInfo(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("获取用户信息失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, co.Success("获取用户信息成功", data))
}

func UpdateUserInfo(c *gin.Context) {
	userId := app.GetUserId(c)
	var data dto.UpdateInfoMap
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	if err := user.UpdateUserInfo(c, userId, data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("用户信息修改失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, co.Success("用户信息修改成功", nil))
}

func UpdatePassword(c *gin.Context) {
	userId := app.GetUserId(c)
	var param dto.UpdatePasswordMap
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	if err := user.UpdatePassword(c, userId, param.Password); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("密码修改失败"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, co.Success("密码修改成功", nil))
}
