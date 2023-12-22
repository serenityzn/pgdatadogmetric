package main

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type PosetgresInterface interface {
	GetConnectionsCount() (int, error)
	GetCount() (int, error)
}

type PostgresService struct {
	db *sql.DB
}

func NewPostgresService(db *sql.DB) *PostgresService {
	return &PostgresService{db: db}
}

func (ps *PostgresService) GetCount() (int, error) {
	if ps.db.Ping() != nil {
		sysStat.DbStatus = "no connection"
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error("Connection to Database Lost")
		return 0, nil
	}
	sysStat.DbStatus = "connected"
	var count int

	row := ps.db.QueryRow("select count(*) from tpart_charge_20231220w3 where end_time > '2023-12-20 14:28:05.986'")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ps *PostgresService) GetConnectionsCount() (int, error) {
	if ps.db.Ping() != nil {
		return 0, nil
	}
	var count int
	row := ps.db.QueryRow("select count(*) from pg_stat_activity")
	err := row.Scan(&count)
	if err != nil {
		log.WithFields(log.Fields{
			"LogLevel": "error",
		}).Error(err)

		return 0, err
	}
	return count, nil
}
