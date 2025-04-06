package dao

import "time"

type WritingReports struct {
	ID            int       `json:"id" gorm:"column:id;primaryKey"`
	UserID        int64     `json:"userId" gorm:"column:user_id;not null"`
	ReportContent string    `json:"reportContent" gorm:"column:report_content"`
	ReportTime    time.Time `json:"reportTime" gorm:"column:report_time;not null"`
}

func (WritingReports) TableName() string {
	return "writing_reports"
}
