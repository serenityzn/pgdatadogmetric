package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var logLevelMap = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"panic": log.PanicLevel,
}

func logInit() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	sysStat.LogLevel = "debug"

	logfile, err := os.OpenFile("./application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.SetOutput(os.Stdout)
		logWF("error", err.Error(), "logging.logInit")
	} else {
		log.SetOutput(logfile)
		//logfile.Close()
		//log.SetOutput(os.Stdout)
	}
	logWF("info", "Logging Initialized", "logging.logInit")
}
