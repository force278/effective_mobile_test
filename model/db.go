package model

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

type DB struct {
	Conn *pgx.Conn
}

func (d *DB) Init() error {
	// Установка соединения с БД
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDATABASE"))
	connConfig, err := pgx.ParseConnectionString(databaseUrl)
	if err != nil {
		return err
	}
	d.Conn, err = pgx.Connect(connConfig)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) Close() error {
	if d.Conn != nil {
		return d.Conn.Close()
	}
	return nil
}
