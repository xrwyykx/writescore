package global

import (
	"context"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func GetDbConn(c context.Context) *gorm.DB {
	return dbConn.WithContext(c)
}
