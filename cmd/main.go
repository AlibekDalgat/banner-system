package main

import (
	"banner-system/internal/cache"
	"banner-system/internal/config"
	"banner-system/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка при получении переменных окружения %s", err.Error())
	}
	dbCfg := config.GetDBConfig()
	db, err := repository.OpenDB(dbCfg)
	if err != nil {
		logrus.Fatalf("Ошибка при подклюении к базе данных: %s", err.Error())
	}
	cacheOpts, err := config.GetCacheOpts()
	if err != nil {
		logrus.Fatalf("Ошибка в переменных окружения для подключения к кэшу: %s", err.Error())
	}
	cache, err := cache.OpenCache(cacheOpts)
	if err != nil {
		logrus.Fatalf("Ошибка при подклюении к кэшу: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := delivery.NewHandler(services, cache)
	_ = handlers
}
