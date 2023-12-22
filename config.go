package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func configInit() (pgConnect, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./etc/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)
		return pgConnect{}, err
	}

	var config pgConnect

	viper.BindEnv("pg.host", "PG_HOST")
	viper.BindEnv("pg.port", "PG_PORT")
	viper.BindEnv("pg.dbname", "PG_DBNAME")
	viper.BindEnv("pg.user", "PG_USER")
	viper.BindEnv("pg.password", "PG_PASSWORD")

	config.host = viper.GetString("pg.host")
	config.port = viper.GetInt("pg.port")
	config.dbname = viper.GetString("pg.dbname")
	config.user = viper.GetString("pg.user")
	config.password = viper.GetString("pg.password")

	return config, nil
}
