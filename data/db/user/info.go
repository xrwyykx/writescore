package user

import (
	"github.com/gin-gonic/gin"
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
