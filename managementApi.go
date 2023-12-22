package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

var myApp dbApp

func startRouter(app dbApp) error {
	myApp = app
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/v1/version", getVersion)
	r.Get("/v1/health", getHealth)
	r.Get("/v1/mgmt/connections", getConnections)
	r.Get("/v1/mgmt/count", getCount)

	r.Put("/v1/mgmt/loglevel/{level}", setLogLevel)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}

	return nil
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("version 1.0.0"))
	log.WithFields(log.Fields{
		"LogLevel": "debug",
	}).Debug("Version Requested")
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(sysStat)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	log.WithFields(log.Fields{
		"LogLevel": "debug",
	}).Debug("Health Requested")
}

func getConnections(w http.ResponseWriter, r *http.Request) {
	cn, err := myApp.pg.GetConnectionsCount()
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	w.Write([]byte(fmt.Sprintf("Connections %d", cn)))
	log.WithFields(log.Fields{
		"LogLevel": "debug",
	}).Debug("Connections Requested")
}

func getCount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Connections %d", count)))
	log.WithFields(log.Fields{
		"LogLevel": "debug",
	}).Debug("Count Requested")
}

func setLogLevel(w http.ResponseWriter, r *http.Request) {
	level := chi.URLParam(r, "level")
	log.SetLevel(logLevelMap[level])
	sysStat.LogLevel = level
	w.Write([]byte(fmt.Sprintf("Log Level set to %s", level)))
	log.WithFields(log.Fields{
		"LogLevel": "debug",
	}).Debug("Log Level Set")
}
