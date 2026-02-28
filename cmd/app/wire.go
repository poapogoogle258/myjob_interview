//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/poapogoogle258/myjob_interview/internel/handler"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
)

func initializeServer(db *mongo.Database) *gin.Engine {
	wire.Build(
		repository.NewJobRepository,
		handler.NewJobHandler,
		NewRouter,
	)
	return nil
}
