package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type dbCbreakerStruct struct {
	alive        bool
	retryTimeout int
	failCount    int
}

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
var exit = 0
var sysStat = systemConfig{}

var br = dbCbreakerStruct{
	alive:        false,
	retryTimeout: 1,
	failCount:    0,
}

var myServerContext context.Context
var myServerClose context.CancelFunc

const waitTime = 5

func main() {
	logInit()

	myDatabaseConf, err := configInit()
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)
		return
	}

	db, err := dbConnect(myDatabaseConf)
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)
		return
	} else {
		sysStat.DbStatus = "connected"
		logWF("info",
			fmt.Sprintf("Connected to Database [%s] on Host [%s]", myDatabaseConf.dbname, myDatabaseConf.host),
			"main")
	}

	defer db.Close()

	mainContext, mainContextCancel := context.WithCancel(context.Background())
	defer mainContextCancel()

	mainContext = context.WithValue(mainContext, "Config", myDatabaseConf)

	go dbCircuitBreaker(mainContext, db)

	application := dbApp{pg: NewPostgresService(db)}

	myServerContext, myServerClose = context.WithCancel(mainContext)

	httpSrv := managementService{
		server: http.Server{},
	}
	httpSrv.serverInit()

	go httpSrv.startRouter(application)

	go dbRoutine(application)

	go dbPing(application)

	go func() { //TODO READ CHANNELS

		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-myServerContext.Done():
				return
			case <-ticker.C:
				if exit == 1 {

					err = httpSrv.server.Shutdown(myServerContext)
					if err != nil {
						logWF("error", err.Error(), "main.shutdownTicker")
					}
					myServerClose()
				}
			default:
			}
		}

	}()

	for {
		if exit == 1 {
			mainContextCancel()
			time.Sleep(time.Second * 1)
			break
		}
		time.Sleep(time.Second * 5)
		//err = datadogSubmitMetric(float64(count))
		//if err != nil {
		//		panic(err)
		//	}
		logWF("info", fmt.Sprintf("Count: %d", count), "main")
	}

}

func dbConnect(connect pgConnect) (*sql.DB, error) {

	psqlConnect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connect.host, connect.port, connect.user, connect.password, connect.dbname)

	db, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logWF("error", err.Error(), "main.dbConnect")
		return nil, err
	}
	logWF("debug",
		fmt.Sprintf("Connected to Database [%s] on Host [%s]", connect.dbname, connect.host),
		"main.dbConnect")
	return db, nil
}

func dbRoutine(app dbApp) {
	for {

		countResult, err := app.pg.GetCount()
		if err != nil {
			logWF("error", err.Error(), "main.dbRoutine")
		}
		count = countResult
		time.Sleep(time.Second * waitTime)
	}
}

func dbPing(app dbApp) {
	for {
		sysStat.DbStatus = "Connected."
		err := app.pg.DbPing()
		if err != nil {
			sysStat.DbStatus = "Failed to connect"
		}
		time.Sleep(time.Second * 2)
	}
}

func dbCircuitBreaker(ctx context.Context, db *sql.DB) {
	for {
		select {
		case <-ctx.Done():
			logWF("error", ctx.Err().Error(), "main.dbCircuitBreaker")
			return
		default:
			for {
				err := db.Ping()
				if err == nil {
					br.alive = true
					br.failCount = 0
					break
				}
				br.alive = false
				br.failCount += 1
				time.Sleep(time.Second * 1)
				if br.failCount > 3 {
					br.failCount = 0
					br.retryTimeout += 2
					break
				}
			}
			time.Sleep(time.Second * time.Duration(br.retryTimeout))
			//TODO Rewriter using ticker

		}

	}

}

func logWF(logLevel string, message string, system string) {
	switch logLevel {
	case "info":
		log.WithFields(log.Fields{
			"LogLevel": logLevel,
			"System":   system,
		}).Info(message)
	case "error":
		log.WithFields(log.Fields{
			"LogLevel": logLevel,
			"System":   system,
		}).Error(message)
	case "debug":
		log.WithFields(log.Fields{
			"LogLevel": logLevel,
			"System":   system,
		}).Debug(message)
	}
}
