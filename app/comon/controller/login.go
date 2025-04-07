package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"writescore/data/db/comon"
	"writescore/models/co"
	"writescore/models/dto"
)

func Register(c *gin.Context) {
	var data dto.RegisterMap
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	log.Printf("Received register data: %+v\n", data)
	if err := comon.Register(c, data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("用户注册失败"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, co.Success("用户注册成功", nil))
}
func Login(c *gin.Context) {
	var data dto.LoginMap
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("参数绑定失败"+err.Error()))
		return
	}
	if err := comon.CheckLogin(c, data); err != nil {
		c.JSON(http.StatusBadRequest, co.BadRequest("用户登录失败"+err.Error()))
		return
	}
	//token := c.MustGet("token").(string)
	c.JSON(http.StatusOK, co.Success("登录成功", gin.H{
		//"token": token,
	}))
}
