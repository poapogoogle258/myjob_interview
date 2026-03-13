package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
	"github.com/poapogoogle258/myjob_interview/internel/usecase"
)

type JobHandler struct {
	usecase *usecase.ScraperUsecase
	repo    repository.JobRepository
}

func NewJobHandler(repo repository.JobRepository, usecase *usecase.ScraperUsecase) *JobHandler {
	return &JobHandler{repo: repo, usecase: usecase}
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {
	jobs, err := h.repo.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	jobID := c.Param("id")

	var body struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !h.repo.IsExist(c.Request.Context(), jobID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	err := h.repo.UpdateStatus(c.Request.Context(), jobID, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "job status updated successfully"})
}

func (h *JobHandler) GetScrapingJobStatus(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		Message      string     `json:"message"`
		IsProcessing bool       `json:"processing"`
		Time         *time.Time `json:"time"`
	}{
		Message:      "success",
		IsProcessing: h.usecase.IsProcessing(),
		Time:         h.usecase.GetScrapingJobLastTime(),
	})
}

func (h *JobHandler) ActiveScrapingManual(c *gin.Context) {

	if h.usecase.IsProcessing() {
		c.JSON(http.StatusOK, gin.H{"message": "scraping is already in progress"})
		return
	}

	go h.usecase.ScrapingJob()

	c.JSON(http.StatusOK, gin.H{"message": "scraping is processing"})
}
