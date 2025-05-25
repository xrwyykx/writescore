package comon

import (
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
	"writescore/global"
	"writescore/models/dao"
	"writescore/models/dto"
)

func Register(c *gin.Context, param dto.RegisterMap) error {
	//给用户的密码进行加密，并且初始化时间，通过雪花算法生成id
	if param.Username == "" || param.Password == "" {
		return errors.New("请输入正确的用户名或密码")
	}
	if param.Avatar == "" {
		return errors.New("速速上传一个可爱的头像")
	}
	var count int64
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("username = ?", param.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名重复咯，快快换一个！！")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userId, err := GenerateSnowflakeId()
	if err != nil {
		return err
	}
	user := dao.User{
		ID:         userId,
		Username:   param.Username,
		Password:   string(hashedPassword),
		CreateTime: time.Now(),
		NickName:   param.NickName,
		Avatar:     param.Avatar,
	}
	if err := global.GetDbConn(c).Model(&dao.User{}).Create(&user).Error; err != nil {
		return err
	}
	return nil
}
func GenerateSnowflakeId() (int64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0, err
	}
	return node.Generate().Int64(), nil
}

func CheckLogin(c *gin.Context, param dto.LoginMap) error {
	var data dao.User
	if param.Username == "" || param.Password == "" {
		return errors.New("请输入正确的账号或者密码")
	}
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("username = ?", param.Username).First(&data).Error; err != nil {
		return err //找不到的时候也会返回错误
	}
	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(param.Password)); err != nil {
		return errors.New("密码不对奥！")
	}

	//如果验证通过，那么接下来就要给他生成jwt了
	//token, err := middlewares.GEnerateJWT(data.ID)
	//session := generateSession()
	//saveSessionToRedis(c, session, data.Username)
	//c.SetCookie("SESSION", session, 3600, "/", "", false, true)

	return nil
}

//func generateSession() string {
//	return uuid.New().String()
//}
//func saveSessionToRedis(c *gin.Context, session string, userName string) {
//	redisClient := global.GetRedisConn()
//	var user dao.User
//	if err := global.GetDbConn(c).Model(&dao.User{}).Where("username = ?", userName).First(&user).Error; err != nil {
//		log.Println("获取用户信息失败")
//		return
//	}
//	userJson, err := json.Marshal(user)
//	if err != nil {
//		log.Println("序列化用户信息失败:", err)
//		return
//	}
//	//err = redisClient.HSet(c, global.ProjectName+"sessions:"+session, "sessionAttr:user_login", string(userJson)).Err()
//	err = redisClient.HSet(c, global.ProjectName+":sessions:"+session, "sessionAttr:user_login", string(userJson)).Err()
//	if err != nil {
//		log.Println("将sessioin存入redis失败:", err)
//	}
//}
