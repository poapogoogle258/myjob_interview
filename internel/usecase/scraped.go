package usecase

import (
	"encoding/json"
	"fmt"

	_ "github.com/poapogoogle258/myjob_interview/internel/clients/jobsdb"
	_ "github.com/poapogoogle258/myjob_interview/internel/clients/jobthai"
	"github.com/poapogoogle258/myjob_interview/internel/model"
	"github.com/poapogoogle258/myjob_interview/internel/provider/clients_provider"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
	"github.com/poapogoogle258/myjob_interview/internel/tool/ollama"
)

type ScraperUsecase struct {
	repo repository.JobRepository
}

func NewScraperUsecase(repo repository.JobRepository) *ScraperUsecase {
	return &ScraperUsecase{
		repo: repo,
	}
}

func analysisSkill(jobDescription string) (*model.SkillsModel, error) {
	jsonString, err := ollama.GetSkillsRequestFromContent(jobDescription)
	if err != nil {
		return nil, err
	}

	skills := &model.SkillsModel{}
	if err = json.Unmarshal([]byte(jsonString), skills); err != nil {
		return nil, err
	}

	return skills, nil
}

func (u *ScraperUsecase) ScrapingJob() {
	fmt.Println("scraping ..")

	for _, client := range clients_provider.GetListProvider() {
		provider, ok := clients_provider.GetProvider(client)
		if !ok {
			continue
		}

		fmt.Printf("fetch data %s ", provider.GetName())
		jobs, err := provider.FetchJobs()
		if err != nil {
			fmt.Printf("\nError: Client %s error: %s\n", provider.GetName(), err)
			continue
		}
		fmt.Printf("found %d jobs.\n", len(jobs))
		for i, job := range jobs {
			fmt.Printf("analysis(%d/%d) title: %s url: %s\n", i+1, len(jobs), job.Title, job.URL)
			result, err := analysisSkill(job.Description)
			if err != nil {
				fmt.Printf("\nError: Analysis JobId %s Client %s error: %s\n", job.ExternalID, provider.GetName(), err)
				continue
			}
			job.Skills = result

		}

	}
}
