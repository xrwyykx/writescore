package dao

import "time"

type Essay struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey"`
	UserID     int64     `json:"userId" gorm:"column:user_id;not null"`
	Title      string    `json:"title" gorm:"column:title;null"`
	Content    string    `json:"content" gorm:"column:content"`
	SubmitTime time.Time `json:"submitTime" gorm:"column:submit_time;not null"`
	Score      *float64  `json:"score" gorm:"column:score"`
	Feedback   *string   `json:"feedback" gorm:"column:feedback"`
	//WordCount    int       `json:"wordCount" gorm:"column:word_count"`
	//UploadMethod int `json:"uploadMethod" gorm:"column:upload_method"`
	ImageId int `json:"imageId" gorm:"column:image_id"`
}

func (Essay) TableName() string {
	return "essay"
}

//也就是说这里会有一个最终的评分还有综合的反馈，然后在essay_scoring_details表中会有对于每个标准具体的评分，每项标准具体的反馈
