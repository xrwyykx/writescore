package dto

type RegisterMap struct {
	Username string `json:"username" gorm:"column:username;not null"`
	Password string `json:"password" gorm:"column:password;not null"`
	NickName string `json:"nick_name" gorm:"column:nickname"`
	//Avatar   string `json:"avatar" gorm:"column:avatar"`
}

type LoginMap struct {
	Username string `json:"username" gorm:"column:username;not null"`
	Password string `json:"password" gorm:"column:password;not null"`
}
