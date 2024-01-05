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

type managementService struct {
	server http.Server
}

type databaseService dbApp

func (mgmtSrv *managementService) serverInit() {
	mgmtSrv.server.Addr = ":8080"
	mgmtSrv.server.ReadTimeout = time.Second * 10
}

func (mgmtSrv *managementService) startRouter(app dbApp) {

	db := databaseService(app)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/v1/version", getVersion)
	r.Get("/v1/health", getHealth)
	r.Get("/v1/mgmt/connections", db.getConnections)
	r.Get("/v1/mgmt/count", getCount)
	r.Post("/v1/mgmt/exit", programExit)
	r.Put("/v1/mgmt/setlog/{level}", setLogLevel)

	mgmtSrv.server.Handler = r
	err := mgmtSrv.server.ListenAndServe()
	if err != nil {
		return
	}
}

func programExit(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Exiting"))
	exit = 1
	logWF("info", "Exiting", "managementApi.programExit")
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("version 1.0.0"))
	logWF("debug", "Version Requested", "managementApi.getVersion")
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(sysStat)
	if err != nil {
		_, _ = w.Write([]byte("Error: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if sysStat.DbStatus == "no connection" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_, _ = w.Write(response)
	logWF("debug", "Health Requested", "managementApi.getHealth")
}

func (myApp *databaseService) getConnections(w http.ResponseWriter, r *http.Request) {
	cn, err := myApp.pg.GetConnectionsCount()
	if err != nil {
		_, _ = w.Write([]byte("Error: " + err.Error()))
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf("Connections %d", cn)))
	logWF("debug", "Connections Requested.", "managementApi.getConnections")
}

func getCount(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(fmt.Sprintf("Count %d", count)))
	logWF("debug", "Count Requested", "managementApi.getCount")
}

func setLogLevel(w http.ResponseWriter, r *http.Request) {
	level := chi.URLParam(r, "level")
	log.SetLevel(logLevelMap[level])
	sysStat.LogLevel = level
	_, _ = w.Write([]byte(fmt.Sprintf("Log Level set to %s", level)))
	logWF("debug", "Log Level Set", "managementApi.setLogLevel")
}
