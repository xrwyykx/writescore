package dao

import "time"

type EssayScoringDetails struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey"`
	EssayID    int       `json:"essayId" gorm:"column:essay_id;not null"`
	CriteriaID int       `json:"criteriaId" gorm:"column:criteria_id;not null"`
	Score      float64   `json:"score" gorm:"column:score"`
	Feedback   string    `json:"feedback" gorm:"column:feedback"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;not null"`
}

func (EssayScoringDetails) TableName() string {
	return "essay_scoring_details"
}
