package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/markraiter/simple-blog/cmd/server"
	"github.com/markraiter/simple-blog/pkg/handler"
	"github.com/markraiter/simple-blog/pkg/repository"
	"github.com/markraiter/simple-blog/pkg/service"
	"github.com/spf13/viper"
)

// @title REST API for Simple Blog Swagger Example
// @version 1.0
// @description This is a simple blog project for practicing Go REST API.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func Start() {

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("failed to load environment variables: %s/n", err.Error())
	}

	// Initializing Database
	db, err := repository.NewPostgresDB(repository.Config{
		Driver:   viper.GetString("db.driver"),
		Username: viper.GetString("db.username"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASS"),
	})

	if err != nil {
		log.Fatalf("error initializing database: %s\n", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(server.Server)
	go func() {
		if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error ocured while running http server: %s\n", err.Error())
		}
	}()

	log.Print("Blog Started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Blog Shutting Down")

	if err := server.Shotdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
