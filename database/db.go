package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"together-ai-assistant/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Database,
	)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(time.Hour)
			if err := db.Ping(); err != nil {
				log.Printf("Database ping failed: %v", err)
			}
		}
	}()
	return db.Ping()
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	return db.Close()
}
