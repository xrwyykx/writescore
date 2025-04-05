package global

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DB struct {
		Host     string
		Username string
		Password string
		Port     string
		DBName   string
		Charset  string
	}
	HTTP struct {
		Path string
		Port string
	}
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config", err)
		return
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("failed to unmarshal config", err)
		return
	}
	log.Println("config loaded successfully")
}
