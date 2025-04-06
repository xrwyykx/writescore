package dao

import "time"

type EssayErrors struct {
	ID            int       `json:"id" gorm:"column:id;primaryKey"`
	EssayID       int       `json:"essayId" gorm:"column:essay_id;not null"`
	ErrorType     int       `json:"errorType" gorm:"column:error_type;not null"`
	ErrorContent  string    `json:"errorContent" gorm:"column:error_content"`
	Suggestion    string    `json:"suggestion" gorm:"column:suggestion"`
	ErrorPosition int       `json:"errorPosition" gorm:"column:error_position"`
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time;not null"`
}

func (EssayErrors) TableName() string {
	return "essay_errors"
}
