package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"writescore/global"
	"writescore/models/dao"
	"writescore/models/dto"
	"writescore/utils"
)

func GetUserInfo(c *gin.Context, id int64) (data dto.UserInfo, err error) {
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("id = ?", id).Select("*").First(&data).Error; err != nil {
		return dto.UserInfo{}, err
	}
	data.CreateTimeMar = utils.MarshalTime(data.CreateTime)
	return data, nil
}

func UpdateUserInfo(c *gin.Context, id int64, param dto.UpdateInfoMap) error {
	var count int64

	//看看这个人存不存在
	if param.Username != "" {
		if err := global.GetDbConn(c).Model(&dao.User{}).Where("username = ?", param.Username).Count(&count).Error; err != nil {
			return err
		}
	}
	if count > 0 {
		return errors.New("用户名重复，请重新想一个")
	}
	updateMap := map[string]interface{}{}
	if param.Avatar != "" {
		updateMap["avatar"] = param.Avatar
	}
	if param.Username != "" {
		updateMap["username"] = param.Username
	}
	if param.NickName != "" {
		updateMap["nick_name"] = param.NickName
	}
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("id = ?", id).Updates(&updateMap).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePassword(c *gin.Context, id int64, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err = global.GetDbConn(c).Model(&dao.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password": string(hashedPass),
	}).Error; err != nil {
		return err
	}
	return nil
}
