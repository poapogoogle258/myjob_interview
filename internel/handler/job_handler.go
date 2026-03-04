package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
)

type JobHandler struct {
	repo repository.JobRepository
}

func NewJobHandler(repo repository.JobRepository) *JobHandler {
	return &JobHandler{repo: repo}
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

	err := h.repo.UpdateStatus(c.Request.Context(), jobID, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "job status updated successfully"})
}
