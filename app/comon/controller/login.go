package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"writescore/data/db/comon"
	"writescore/global"
	"writescore/models/co"
	"writescore/models/dao"
	"writescore/models/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	} else {
		session := generateSession()
		saveSessionToRedis(c, session, data.Username)

		// 设置cookie，不指定 domain，让浏览器自动处理
		c.SetCookie("SESSION", session, 3600, "/", "", false, false)

		// 打印调试信息
		log.Printf("Setting cookie: SESSION=%s", session)
		log.Printf("Response headers: %v", c.Writer.Header())
		log.Printf("Request Origin: %s", c.Request.Header.Get("Origin"))

		c.JSON(http.StatusOK, co.Success("登录成功", gin.H{"session": session}))
	}
}

func generateSession() string {
	return uuid.New().String()
}

func saveSessionToRedis(c *gin.Context, session string, userName string) {
	redisClient := global.GetRedisConn()
	var user dao.User
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("username = ?", userName).First(&user).Error; err != nil {
		log.Println("获取用户信息失败")
		return
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		log.Println("序列化用户信息失败:", err)
		return
	}
	err = redisClient.HSet(c, global.ProjectName+":sessions:"+session, "sessionAttr:user_login", string(userJson)).Err()
	if err != nil {
		log.Println("将sessioin存入redis失败:", err)
	}
}
