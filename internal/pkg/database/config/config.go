package config

import (
	"os"
	"strconv"
)

type PsqlDBConnection struct {
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
}

type RedisConnection struct {
	DBAddress  string
	DBPassword string
	DB         int
}

type DatabaseConfig struct {
	Psql  PsqlDBConnection
	Redis RedisConnection
}

type Database DatabaseConfig

func NewDatabase() *DatabaseConfig {
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &DatabaseConfig{
		Psql: PsqlDBConnection{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBDatabase: os.Getenv("DB_NAME"),
			DBUsername: os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
		},
		Redis: RedisConnection{
			DBAddress:  os.Getenv("REDIS_DB_ADDR"),
			DBPassword: os.Getenv("REDIS_DB_PASSWORD"),
			DB:         redisDB,
		},
	}
}
