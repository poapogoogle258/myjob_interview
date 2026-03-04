package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/poapogoogle258/myjob_interview/internel/db"
	"github.com/poapogoogle258/myjob_interview/internel/handler"
	"github.com/robfig/cron/v3"
)

func NewRouter(h *handler.JobHandler) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	v1 := r.Group("/api/v1")
	{
		v1.GET("/job", h.GetAllJobs)
		v1.PUT("/job/:id/status", h.UpdateJobStatus)
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

	app := initializeServer(db)

	c := cron.New()
	c.AddFunc("@every 30m", app.Scraper.ScrapingJob)
	c.Start()
	defer c.Stop()

	address := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")

	if err := app.Router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
