package dto

import (
	"time"
	"writescore/models"
)

type RatingEssayMap struct {
	EssayId int    `json:"essayId" gorm:"column:essay_id"`
	Title   string `json:"title" gorm:"column:title"`
}

type RatingResult struct {
	EssayId       int `json:"essayId" gorm:"column:essay_id"`
	PerScore      []PerScore
	FinalScore    float64 `json:"finalScore" gorm:"column:final_score"`
	FinalFeekback string  `json:"finalFeekback" gorm:"column:final_feekback"`
}

type PerScore struct {
	CriteriaId    int     `json:"criteriaId" gorm:"column:criteria_id"`
	CriteriaName  string  `json:"criteriaName" gorm:"column:criteria_name"`
	CriteriaScore float64 `json:"criteriaScore" gorm:"column:criteria_score"`
	Feekback      string  `json:"feedback" gorm:"column:feedback"`
}

type GetEssayMap struct {
	//按找提交时间起止来寻找，还有总评分范围，文章标题
	StartTime   string  `json:"startTime" gorm:"column:start_time"`
	EndTime     string  `json:"endTime" gorm:"column:end_time"`
	MinScore    float64 `json:"minScore" gorm:"column:min_score"`
	MaxScore    float64 `json:"maxScore" gorm:"column:max_score"`
	models.Page         //实现分页管理
	Title       string  `json:"title" gorm:"column:title"`
}

// 获取到的是所有文章的大致信息就可以了
type AllEssays struct {
	SubmitTime    time.Time `json:"-" gorm:"column:submit_time"`
	SubmitTimeMar string    `json:"submitTime" gorm:"column:submit_time_mar"`
	Id            int       `json:"id" gorm:"column:id"` //文章id
	Score         float64   `json:"score" gorm:"column:score"`
	Title         string    `json:"title" gorm:"column:title"`
}

type EssayDetail struct {
	PerScore []PerScore
	Content  string `json:"content" gorm:"column:content"`
	//SubmitTime    time.Time `json:"-" gorm:"column:submit_time"`
	SubmitTimeMar string   `json:"submitTime" gorm:"column:submitTime"`
	Score         *float64 `json:"score" gorm:"column:score"`
	Feedback      *string  `json:"feedback" gorm:"column:feedback"`
	Title         string   `json:"title" gorm:"column:title"`
}
