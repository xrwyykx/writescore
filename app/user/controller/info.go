package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"writescore/app"
	"writescore/data/db/user"
	"writescore/models/co"
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
