package config

import (
	"log"
	"resumes/internal/database"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("configs/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
}

func GetDBConfig() *database.Config {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	return &database.Config{
		Addr:     viper.GetString("mysql_host"),
		Port:     viper.GetUint16("mysql_port"),
		Password: viper.GetString("mysql_password"),
		User:     viper.GetString("mysql_user"),
		DB:       viper.GetString("mysql_db"),
	}
}
