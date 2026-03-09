package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
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
	processing bool
	lastTime   *time.Time
	logger     *slog.Logger
	repo       repository.JobRepository
}

func NewScraperUsecase(repo repository.JobRepository, logger *slog.Logger) *ScraperUsecase {
	now := time.Now()
	return &ScraperUsecase{
		processing: false,
		lastTime:   &now,
		repo:       repo,
		logger:     logger,
	}
}

func (u *ScraperUsecase) IsProcessing() bool {
	return u.processing
}

func (u *ScraperUsecase) GetScrapingJobLastTime() *time.Time {
	return u.lastTime
}

func (u *ScraperUsecase) ScrapingJob() {
	if u.processing {
		return
	}

	u.processing = true
	defer func() {
		u.processing = false
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	u.logger.InfoContext(ctx, "start scraping schedule job")
	for _, client := range clients_provider.GetListProvider() {
		provider, ok := clients_provider.GetProvider(client)
		if !ok {
			continue
		}
		u.logger.InfoContext(ctx, fmt.Sprintf("start scraping client %s", provider.GetName()))
		jobs, err := provider.FetchJobs()
		if err != nil {
			u.logger.WarnContext(ctx, fmt.Sprintf("scraping client %s failed", provider.GetName()), "error", err)
			continue
		}

		u.logger.InfoContext(ctx, fmt.Sprintf("website %s found %d jobs.", provider.GetName(), len(jobs)))
		length_job := len(jobs)

		for i, job := range jobs {
			jobId := model.GetHashJobId(job)
			job.HashId = jobId
			result, _ := u.repo.GetByHashId(ctx, jobId)
			if result == nil {
				u.logger.InfoContext(ctx, fmt.Sprintf("(%d/%d)analysis job-id %s", i+1, len(jobs), jobId), "client", provider.GetName(), "externalID", job.ExternalID, "title", job.Title, "company", job.CompanyName)
				skills, err := AnalysisSkill(job.Description)
				if err != nil {
					u.logger.WarnContext(ctx, fmt.Sprintf("(%d/%d)analysis job-id %s failed", i+1, len(jobs), jobId), "error", err, "client", provider.GetName(), "externalID", job.ExternalID, "title", job.Title, "company", job.CompanyName)
					continue
				}
				job.Status = "new"
				job.Skills = skills
			} else {
				job.Status = result.Status
				job.Skills = result.Skills
			}

			if !slices.Contains(job.Skills.Languages, "golang") {
				length_job--
			}

			u.repo.UpsertByExternalID(ctx, job)
			u.logger.InfoContext(ctx, fmt.Sprintf("(%d/%d)job-id %s updated", i+1, len(jobs), jobId))
		}

		u.logger.InfoContext(ctx, fmt.Sprintf("website %s found %d jobs have golang in attribute and updated", provider.GetName(), length_job))
		now := time.Now()
		u.lastTime = &now
	}
}

func AnalysisSkill(jobDescription string) (*model.SkillsModel, error) {
	prompt := fmt.Sprintf(`
	Summarize the required skills from this job description, listing them item by item order by priority using the json format (ToLower)
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
	ctx_10s, cancel_10s := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel_10s()
	if err := client.Generate(ctx_10s, req, respFunc); err != nil {
		return nil, err
	}

	skills := &model.SkillsModel{}
	if err = json.Unmarshal([]byte(responseBuilder.String()), skills); err != nil {
		return nil, err
	}

	return skills, nil
}
