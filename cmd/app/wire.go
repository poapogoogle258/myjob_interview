//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/poapogoogle258/myjob_interview/internel/handler"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
	"github.com/poapogoogle258/myjob_interview/internel/usecase"
)

type App struct {
	Router  *gin.Engine
	Scraper *usecase.ScraperUsecase
}

func initializeServer(db *mongo.Database) *App {
	wire.Build(
		repository.NewJobRepository,
		handler.NewJobHandler,
		usecase.NewScraperUsecase,
		NewRouter,
		wire.Struct(new(App), "*"),
	)
	return nil
}
