package dao

type Administrators struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey"`
	Username string `json:"username" gorm:"column:username;not null"`
	Password string `json:"password" gorm:"column:password;not null"`
}

func (Administrators) TableName() string {
	return "administrators"
}
