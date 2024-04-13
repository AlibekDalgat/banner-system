package main

import (
	"banner-system/internal/app"
	"banner-system/internal/cache"
	"banner-system/internal/config"
	delivery "banner-system/internal/handler"
	"banner-system/internal/repository"
	"banner-system/internal/service"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	srv := new(app.Server)
	go func() {
		if err := srv.Run(os.Getenv("HTTP_PORT"), handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Ошибка при работе http-сервера: %s", err.Error())
		}
	}()

	logrus.Println("Сервис баннеров начал работу")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Маркетплейс завершил работу")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Произошла ошибка при завершении работы сервера: %s", err.Error())
	}
	if status := cache.FlushDB(context.Background()); status.Err() != nil {
		logrus.Errorf("Ошибка при очистке кэша: %s", status.Err().Error())
	}
	if err := cache.Close(); err != nil {
		logrus.Errorf("Ошибка при отсоединении от кэша: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Ошибка при отсоединении от базы данных: %s", err.Error())
	}
}
