package dao

import "time"

type Essay struct {
	ID           int       `json:"id" gorm:"column:id;primaryKey"`
	UserID       int64     `json:"userId" gorm:"column:user_id;not null"`
	Title        string    `json:"title" gorm:"column:title;not null"`
	Content      string    `json:"content" gorm:"column:content"`
	LanguageType int       `json:"languageType" gorm:"column:language_type;not null"`
	SubmitTime   time.Time `json:"submitTime" gorm:"column:submit_time;not null"`
	Score        *float64  `json:"score" gorm:"column:score"`
	Feedback     *string   `json:"feedback" gorm:"column:feedback"`
	OcrResult    *string   `json:"ocrResult" gorm:"column:ocr_result"`
	OcrCorrected bool      `json:"ocrCorrected" gorm:"column:ocr_corrected;default:false"`
	WordCount    int       `json:"wordCount" gorm:"column:word_count"`
}

func (Essay) TableName() string {
	return "essay"
}
