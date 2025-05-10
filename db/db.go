package db

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	// open connection withour fatal
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	// test koneksi untuk memastikan berhasil
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Konfigurasi connection pool (dapat disesuaikan)
	db.SetMaxOpenConns(30)                  // Maximum number of open connections
	db.SetMaxIdleConns(10)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(30 * time.Minute) // Lifetime of each connection

	return db, nil
}
