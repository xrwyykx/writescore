package dto

import "time"

type RestoreImageInfoMap struct {
	ImageURL  string `json:"imageUrl" gorm:"column:image_url;not null"`
	ImageName string `json:"imageName" gorm:"column:image_name;not null"`
	Title     string `json:"title" gorm:"column:title"`
}

type ImageToEssay struct {
	//UserID  int64  `json:"userId" gorm:"column:user_id;not null"`
	EssayId       int       `json:"essayId" gorm:"column:essay_id"`
	Content       string    `json:"content" gorm:"column:content"`
	SubmitTime    time.Time `json:"-" gorm:"column:submit_time;"`
	SubmitTimeMar string    `json:"submitTime"` //图片上传时间
}

type Shangchuan struct {
	ImageURL string `json:"imageUrl" gorm:"column:image_url"`
}
