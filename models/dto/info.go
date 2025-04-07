package dto

import "time"

type UserInfo struct {
	ID            int64     `json:"id" gorm:"column:id;primaryKey"`
	Username      string    `json:"username" gorm:"column:username;not null"`
	Password      string    `json:"password" gorm:"column:password;not null"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;not null"`
	CreateTimeMar string    `json:"createTime"`
	NickName      string    `json:"nickName" gorm:"column:nick_name"`
	Avatar        string    `json:"avatar" gorm:"column:avatar"`
}

type UpdateInfoMap struct {
	Username string `json:"username" gorm:"column:username"`
	NickName string `json:"nickName" gorm:"column:nick_name"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
}

type UpdatePasswordMap struct {
	Password string `json:"password" gorm:"column:password"`
}
