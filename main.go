package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type pgConnect struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

type dbApp struct {
	pg PosetgresInterface
}

type systemConfig struct {
	LogLevel string `json:"logLevel"`
	DbStatus string `json:"dbStatus"`
}

var count int
var sysStat = systemConfig{}

const waitTime = 5

func dbConnect(connect pgConnect) (*sql.DB, error) {

	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connect.host, connect.port, connect.user, connect.password, connect.dbname)

	db, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)
		return nil, err
	}

	log.WithFields(log.Fields{
		"LogLevel": "debug",
	}).Debug("Connected to Database [", connect.dbname, "] on Host [", connect.host, "]")

	return db, nil
}

func dbRoutine(app dbApp) {
	for {

		countResult, err := app.pg.GetCount()
		if err != nil {
			log.WithFields(log.Fields{
				"LogLevel": "error",
			}).Error(err)
		}
		count = countResult
		time.Sleep(time.Second * waitTime)
	}
}

func main() {
	logInit()

	mydb, err := configInit()
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)
		return
	}

	db, err := dbConnect(mydb)
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)
		return
	} else {
		sysStat.DbStatus = "connected"
		log.WithFields(log.Fields{
			"LogLevel": "info",
		}).Info("Connected to Database [", mydb.dbname, "] on Host [", mydb.host, "]")
	}

	defer db.Close()

	application := dbApp{pg: NewPostgresService(db)}

	go startRouter(application)

	go dbRoutine(application)

	for {
		time.Sleep(time.Second * 5)
		//err = datadogSubmitMetric(float64(count))
		//if err != nil {
		//		panic(err)
		//	}

		log.WithFields(log.Fields{
			"LogLevel": "info",
		}).Info("Count: ", count)
	}

}
