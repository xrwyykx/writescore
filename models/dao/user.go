package dao

import "time"

type User struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey"`
	Username   string    `json:"username" gorm:"column:username;not null"`
	Password   string    `json:"password" gorm:"column:password;not null"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;not null"`
	NickName   string    `json:"nickName" gorm:"column:nick_name"`
	Avatar     string    `json:"avatar" gorm:"column:avatar"`
}

func (User) TableName() string {
	return "user"
}
