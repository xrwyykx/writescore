package dao

type ScoringCriteria struct {
	ID           int     `json:"id" gorm:"column:id;primaryKey"`
	CriteriaName int     `json:"criteriaName" gorm:"column:criteria_name;not null"`
	Weight       float64 `json:"weight" gorm:"column:weight;not null"`
}

func (ScoringCriteria) TableName() string {
	return "scoring_criteria"
}
