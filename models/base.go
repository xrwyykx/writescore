package models

type Page struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

type IdMap struct {
	Id int `json:"id" gorm:"column:id"`
}
