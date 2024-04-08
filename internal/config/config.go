package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		User:     os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"),
		Host:     os.Getenv("PGHOST"),
		Port:     os.Getenv("PGPORT"),
		Database: os.Getenv("PGDATABASE"),
		SSLMode:  os.Getenv("PGSSLMODE"),
	}
}

type CacheOpts struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func GetCacheOpts() (CacheOpts, error) {
	db, err := strconv.Atoi(os.Getenv("REDDB"))
	if err != nil {
		return CacheOpts{}, err
	}
	return CacheOpts{
		Host:     os.Getenv("REDHOST"),
		Port:     os.Getenv("REDPORT"),
		Password: os.Getenv("REDPASSWORD"),
		DB:       db,
	}, nil
}
