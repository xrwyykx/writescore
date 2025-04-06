package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"writescore/global"
	"writescore/models/dao"
)

func GetUserId(c *gin.Context) int64 { //返回当前使用者ID
	Session, err := c.Cookie("SESSION")
	if err != nil {
		log.Printf("Failed to get SESSION cookie: %v", err)
		return -1
	}
	redisCil := global.GetRedisConn() //获取redis客户端连接，这个函数返回一个redis客户端实例
	if redisCil == nil {
		log.Println("Redis client is nil")
		return -1
	}
	response := redisCil.HGet(c, global.ProjectName+":sessions:"+Session, "sessionAttr:user_login")
	var userSessionValue dao.User
	err = json.Unmarshal([]byte(response.Val()), &userSessionValue)
	if err != nil {
		return -1
	}
	return userSessionValue.ID
}
