package db

import (
	"database/sql"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DB struct {
	username string
	password string
	host     string
	port     string
	name     string
	Client   *gorm.DB
}

type CloseConn func()

func closeConn(conn *sql.DB) CloseConn {
	//fmt.Println("closeConn")
	return func() {
		_ = conn.Close()
	}
}

func NewDatabase(username, password, host, port, name string) (*gorm.DB, error) {
	logLevel := logger.Silent
	if config.Debug {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, username, password, name, port)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	return client, err
}

func GetConnection() (*gorm.DB, CloseConn) {
	var err error
	db, err := NewDatabase(
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	return db, closeConn(sqlDb)
}
