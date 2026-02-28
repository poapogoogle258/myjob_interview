package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/poapogoogle258/myjob_interview/internel/db"
	"github.com/poapogoogle258/myjob_interview/internel/handler"
)

func NewRouter(h *handler.JobHandler) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/job", h.GetAllJobs)
	}
	return r
}

func main() {

	if err := godotenv.Load(); err != nil { // default load is .env
		panic(err)
	}

	db, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	server := initializeServer(db)
	address := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")

	if err := server.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
