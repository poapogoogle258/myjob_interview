package handler

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
	"github.com/poapogoogle258/myjob_interview/internel/service"
)

type JobHandler struct {
	service *service.ScraperService
	repo    repository.JobRepository
	logger  *slog.Logger
}

func NewJobHandler(repo repository.JobRepository, service *service.ScraperService, logger *slog.Logger) *JobHandler {
	return &JobHandler{repo: repo, service: service, logger: logger}
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
	h.logger.Info("API UpdateJobStatus Route")

	jobID := c.Param("id")

	var body struct {
		Status string `json:"status" binding:"required"`
	}

	h.logger.Info("API UpdateJobStatus Body", "body", body)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status_list := []string{"new", "viewed", "favorite", "registered", "interview", "rejected", "offered", "optional"}
	if !slices.Contains(status_list, body.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	if !h.repo.IsExist(c.Request.Context(), jobID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	err := h.repo.UpdateStatus(c.Request.Context(), jobID, body.Status)
	if err != nil {
		h.logger.Info("API UpdateJobStatus Error", "error", err)
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
		IsProcessing: h.service.IsProcessing(),
		Time:         h.service.GetScrapingJobLastTime(),
	})
}

func (h *JobHandler) ActiveScrapingManual(c *gin.Context) {
	if h.service.IsProcessing() {
		c.JSON(http.StatusOK, gin.H{"message": "scraping is already in progress"})
		return
	}

	go h.service.ScrapingJob()

	c.JSON(http.StatusOK, gin.H{"message": "scraping is processing"})
}
