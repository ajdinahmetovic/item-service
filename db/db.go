package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ajdinahmetovic/item-service/logger"
	//Postgress import
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tajna"
	dbname   = "db_1"
)

var db *sql.DB

//ConnectDB func
func ConnectDB() error {
	psqlInfo := fmt.Sprintf("port=%d user=%s "+
		"password=%s dbname=%s host=%s sslmode=disable",
		port, user, password, dbname, host)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Error("Failed to open connection to database", "time", time.Now(), "err", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		logger.Error("Failed to ping database", "time", time.Now(), "err", err)
		return err
	}
	logger.Info("Database connected", "time", time.Now(), "err", err)
	return nil
}
