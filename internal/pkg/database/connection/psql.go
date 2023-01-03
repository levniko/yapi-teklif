package database

import (
	"fmt"
	"log"
	"time"

	"github.com/yapi-teklif/internal/pkg/database/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *Connection

type IConnection interface {
	PsqlDB() *gorm.DB
}

type Connection struct {
	GetPsqlDB *gorm.DB
}

func (c *Connection) PsqlDB() *gorm.DB {
	return c.GetPsqlDB
}

func Connect() *Connection {
	database_config := config.NewDatabase()
	//gorm_config := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	psql_dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s dbname=%s sslmode=prefer",
		database_config.Psql.DBHost,
		database_config.Psql.DBUsername,
		database_config.Psql.DBPassword,
		database_config.Psql.DBPort,
		database_config.Psql.DBDatabase)

	psql_db, err := gorm.Open(postgres.Open(psql_dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database %v\n, dsn: %s", err, psql_dsn)
	}

	Conn = &Connection{
		GetPsqlDB: psql_db,
	}
	sqlDB, _ := psql_db.DB()
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	fmt.Println("Connected to database")

	return Conn
}
