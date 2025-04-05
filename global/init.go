package global

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		AppConfig.DB.Username,
		AppConfig.DB.Password,
		AppConfig.DB.Host,
		AppConfig.DB.Port,
		AppConfig.DB.DBName,
		AppConfig.DB.Charset,
	)
	var err error
	dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect DB", err)
		return
	}
	DB = dbConn
	log.Println("init db successsfully")
}
