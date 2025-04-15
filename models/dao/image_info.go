package dao

import "time"

type ImageInfo struct {
	ID int `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	//EssayID    int       `json:"essayId" gorm:"column:essay_id;null"`
	ImageURL   string    `json:"imageUrl" gorm:"column:image_url;not null"`
	ImageName  string    `json:"imageName" gorm:"column:image_name;not null"`
	UploadTime time.Time `json:"uploadTime" gorm:"column:upload_time;not null"`
}

func (ImageInfo) TableName() string {
	return "image_info"
}

//流程是照片上传->
