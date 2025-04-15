package dto

type UploadEssayByHandMap struct {
	//UserID       int64  `json:"userId" gorm:"column:user_id;not null"`
	Title        string `json:"title" gorm:"column:title;not null"`
	Content      string `json:"content" gorm:"column:content"`
	LanguageType int    `json:"languageType" gorm:"column:language_type;not null"`
	UploadMethod int    `json:"uploadMethod" gorm:"column:upload_method"`
}
