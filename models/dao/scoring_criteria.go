package dao

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/google/uuid"
	_ "gorm.io/gorm"
)

type ScoringCriteria struct {
	ID           int     `json:"id" gorm:"column:id;primaryKey"`
	CriteriaName string  `json:"criteriaName" gorm:"column:criteria_name;not null;type:ENUM('语法','结构','内容','创意')"` // 修改字段类型为 string，并添加 ENUM 约束
	Weight       float64 `json:"weight" gorm:"column:weight;not null"`
}

func (ScoringCriteria) TableName() string {
	return "scoring_criteria"
}

// 自定义验证函数
func (s *ScoringCriteria) Validate() error {
	validate := validator.New()
	err := validate.Var(s.CriteriaName, "oneof=语法 结构 内容 创意")
	if err != nil {
		return err
	}
	return nil
}
