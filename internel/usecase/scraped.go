package usecase

import (
	"fmt"

	"github.com/poapogoogle258/myjob_interview/internel/repository"
)

type ScraperUsecase struct {
	repo repository.JobRepository
}

func NewScraperUsecase(repo repository.JobRepository) *ScraperUsecase {
	return &ScraperUsecase{
		repo: repo,
	}
}

func (u *ScraperUsecase) ScrapingJob() {
	fmt.Println("scraping ..")
}
