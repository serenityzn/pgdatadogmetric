package main

import (
	"database/sql"
)

type PosetgresInterface interface {
	GetConnectionsCount() (int, error)
	GetCount() (int, error)
	DbPing() error
}

type PostgresService struct {
	db *sql.DB
}

func NewPostgresService(db *sql.DB) *PostgresService {
	return &PostgresService{db: db}
}

func (ps *PostgresService) GetCount() (int, error) {
	var count int

	row := ps.db.QueryRow("select count(*) from tpart_charge_20231220w3 where end_time > '2023-12-20 14:28:05.986'")
	err := row.Scan(&count)
	if err != nil {
		logWF("error", err.Error(), "database.GetCount")
		return 0, err
	}
	return count, nil
}

func (ps *PostgresService) GetConnectionsCount() (int, error) {
	var count int
	row := ps.db.QueryRow("select count(*) from pg_stat_activity")
	err := row.Scan(&count)
	if err != nil {
		logWF("error", err.Error(), "database.GetConnectionsCount")
		return 0, err
	}
	return count, nil
}

func (ps *PostgresService) DbPing() error {
	return ps.db.Ping()
}
