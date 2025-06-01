package user

import (
	"errors"
	"writescore/global"
	"writescore/models/dao"
	"writescore/models/dto"
	"writescore/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserInfo(c *gin.Context, id int64) (data dto.UserInfo, err error) {
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("id = ?", id).Select("*").First(&data).Error; err != nil {
		return dto.UserInfo{}, err
	}
	data.CreateTimeMar = utils.MarshalTime(data.CreateTime)
	return data, nil
}

func UpdateUserInfo(c *gin.Context, id int64, param dto.UpdateInfoMap) (data dto.UserInfo, err error) {
	var count int64

	//看看这个人存不存在
	if param.Username != "" {
		// 检查用户名是否重复，但排除当前用户自己的用户名
		if err := global.GetDbConn(c).Model(&dao.User{}).
			Where("username = ? AND id != ?", param.Username, id).
			Count(&count).Error; err != nil {
			return dto.UserInfo{}, err
		}
	}
	if count > 0 {
		return dto.UserInfo{}, errors.New("用户名重复，请重新想一个")
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
		return dto.UserInfo{}, err
	}

	// 获取更新后的用户信息
	if err := global.GetDbConn(c).Model(&dao.User{}).Where("id = ?", id).Select("*").First(&data).Error; err != nil {
		return dto.UserInfo{}, err
	}
	data.CreateTimeMar = utils.MarshalTime(data.CreateTime)
	return data, nil
}

func UpdatePassword(c *gin.Context, id int64, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//这个接口没有登录但是可以正常使用，只是无法返回错误，会显示修改成功
	if err = global.GetDbConn(c).Model(&dao.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password": string(hashedPass),
	}).Error; err != nil {
		return err
	}
	return nil
}
