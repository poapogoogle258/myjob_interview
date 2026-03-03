package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	_ "github.com/poapogoogle258/myjob_interview/internel/client/jobsdb"
	_ "github.com/poapogoogle258/myjob_interview/internel/client/jobthai"
	"github.com/poapogoogle258/myjob_interview/internel/model"
	"github.com/poapogoogle258/myjob_interview/internel/provider/clients_provider"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
)

type ScraperUsecase struct {
	logger *slog.Logger
	repo   repository.JobRepository
}

func NewScraperUsecase(repo repository.JobRepository, logger *slog.Logger) *ScraperUsecase {
	return &ScraperUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *ScraperUsecase) ScrapingJob() {
	for _, client := range clients_provider.GetListProvider() {
		provider, ok := clients_provider.GetProvider(client)

		if !ok {
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		u.logger.InfoContext(ctx, fmt.Sprintf("start scraping client %s", provider.GetName()))

		jobs, err := provider.FetchJobs()
		if err != nil {
			u.logger.WarnContext(ctx, fmt.Sprintf("scraping client %s failed", provider.GetName()), "error", err)
			continue
		}
		u.logger.InfoContext(ctx, fmt.Sprintf("website %s found %d jobs.", provider.GetName(), len(jobs)))

		for i, job := range jobs {
			jobId := model.GetHashJobId(job)
			if result, _ := u.repo.GetByHashId(ctx, jobId); result != nil {
				job.HashId = jobId
				job.Skills = result.Skills
				u.repo.UpsertByExternalID(ctx, job)
				continue
			}
			ctx_10s, cancel_10s := context.WithTimeout(ctx, 30*time.Second)
			defer cancel_10s()
			u.logger.InfoContext(ctx_10s, fmt.Sprintf("analysis(%d/%d) id: %s title: %s company: %s", i+1, len(jobs), jobId, job.Title, job.CompanyName))
			result, err := AnalysisSkill(ctx_10s, job.Description)
			ctx_10s.Done()
			if err != nil {
				u.logger.WarnContext(ctx_10s, fmt.Sprintf("analysis job-id %s failed", jobId), "error", err, "client", provider.GetName(), "externalID", job.ExternalID, "title", job.Title, "company", job.CompanyName)
				continue
			}

			job.HashId = jobId
			job.Skills = result
			u.repo.UpsertByExternalID(ctx, job)
			u.logger.InfoContext(ctx_10s, fmt.Sprintf("analysis job-id %s success", jobId), "client", provider.GetName(), "externalID", job.ExternalID, "title", job.Title, "company", job.CompanyName)
		}

		ctx.Done()
	}
}

func AnalysisSkill(ctx context.Context, jobDescription string) (*model.SkillsModel, error) {
	prompt := fmt.Sprintf(`
	Summarize the required skills from this job description, listing them item by item order by priority using the json format 
	{ 
		languages  : [ language1, language2, ...],
		frameworks : [ framework1, framework2, ... ],
		tools      : [ tool1, tool2, ... ],
		databases  : [ database1, database2, ... ],
		hardSkills : [ skill1, skill2, ... ],
		softSkills : [ skill1, skill2, ... ]
	} 
	If it has an abbreviation like "mainSkill (subSkill1, subSkill2)," separate it only subSkill into a new skill.
	If it don't have some skill fill [] (list length zero to avoid null)
	job description : 
	%s
	`, jobDescription)

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	steam := false
	req := &api.GenerateRequest{
		Model:  "scb10x/typhoon2.5-qwen3-4b:latest", // Specify the model to use
		Stream: &steam,
		Format: json.RawMessage(`"json"`),
		Prompt: prompt,
	}

	var responseBuilder strings.Builder
	respFunc := func(resp api.GenerateResponse) error {
		responseBuilder.WriteString(resp.Response)
		return nil
	}

	// Use the Chat function to get a response
	if err := client.Generate(ctx, req, respFunc); err != nil {
		return nil, err
	}

	skills := &model.SkillsModel{}
	if err = json.Unmarshal([]byte(responseBuilder.String()), skills); err != nil {
		return nil, err
	}

	return skills, nil
}
