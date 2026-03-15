//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/poapogoogle258/myjob_interview/internel/handler"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
	"github.com/poapogoogle258/myjob_interview/internel/service"

	"github.com/poapogoogle258/myjob_interview/utils/logger"
)

type App struct {
	Router  *gin.Engine
	Scraper *service.ScraperService
}

func initializeServer(db *mongo.Database) *App {
	wire.Build(
		logger.NewLogger,
		repository.NewJobRepository,
		handler.NewJobHandler,
		service.NewScraperUsecase,
		NewRouter,
		wire.Struct(new(App), "*"),
	)
	return nil
}

// go run github.com/google/wire/cmd/wire .
