package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	ollama "github.com/ollama/ollama/api"
	_ "github.com/poapogoogle258/myjob_interview/internel/client/jobsdb"
	_ "github.com/poapogoogle258/myjob_interview/internel/client/jobthai"
	"github.com/poapogoogle258/myjob_interview/internel/client/provider"
	"github.com/poapogoogle258/myjob_interview/internel/model/dao"
	"github.com/poapogoogle258/myjob_interview/internel/repository"
)

type ScraperService struct {
	processing bool
	lastTime   *time.Time
	logger     *slog.Logger
	repo       repository.JobRepository
}

func NewScraperUsecase(repo repository.JobRepository, logger *slog.Logger) *ScraperService {
	now := time.Now()
	return &ScraperService{
		processing: false,
		lastTime:   &now,
		repo:       repo,
		logger:     logger,
	}
}

func (u *ScraperService) IsProcessing() bool {
	return u.processing
}

func (u *ScraperService) GetScrapingJobLastTime() *time.Time {
	return u.lastTime
}

func (u *ScraperService) ScrapingJob() {
	if u.processing {
		return
	}

	u.processing = true
	defer func() {
		u.processing = false
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start_scraping := time.Now()
	token_scraping_job := generateTimeToken(start_scraping)
	u.logger.InfoContext(ctx, "start scraping schedule job")
	for _, client := range provider.GetListProvider() {
		start_scraping_client := time.Now()
		provider, ok := provider.GetProvider(client)
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

		// find and update new job
		for i, job := range jobs {
			jobId := dao.GetHashJobId(job)
			job.HashId = jobId
			job.SyncID = token_scraping_job
			result, _ := u.repo.GetByHashId(ctx, jobId)
			if result == nil {
				u.logger.InfoContext(ctx, fmt.Sprintf("(%d/%d)analysis job-id %s", i+1, len(jobs), jobId), "client", provider.GetName(), "externalID", job.ExternalID, "title", job.Title, "company", job.CompanyName)
				skills, err := AnalysisSkill(job.Description)
				if err != nil {
					u.logger.WarnContext(ctx, fmt.Sprintf("(%d/%d)analysis job-id %s failed", i+1, len(jobs), jobId), "error", err, "client", provider.GetName(), "externalID", job.ExternalID, "title", job.Title, "company", job.CompanyName)
					continue
				}
				job.Skills = skills
				job.StatusLOG = []dao.StatusLog{}
			} else {
				job.Status = result.Status
				job.Skills = result.Skills
				job.StatusLOG = result.StatusLOG
			}

			if !slices.Contains(job.Skills.Languages, "golang") {
				length_job--
			}

			u.repo.UpsertByExternalID(ctx, job)
			u.logger.InfoContext(ctx, fmt.Sprintf("(%d/%d)job-id %s updated", i+1, len(jobs), jobId))
		}

		// update job not found in ScrapingJob
		jobsNotSync, _ := u.repo.GetExceptSyncId(ctx, client, token_scraping_job)
		u.logger.InfoContext(ctx, fmt.Sprintf("website %s update old job not found in scraping %d jobs.", provider.GetName(), len(jobsNotSync)))
		for i, job := range jobsNotSync {
			provider.SyncJobDetail(job)
			jobId := dao.GetHashJobId(job)
			job.HashId = jobId
			job.SyncID = token_scraping_job
			u.repo.UpsertByExternalID(ctx, job)
			u.logger.InfoContext(ctx, fmt.Sprintf("(%d/%d)job-id %s updated", i+1, len(jobsNotSync), jobId))
		}

		time_used := time.Since(start_scraping_client)
		u.logger.InfoContext(ctx, fmt.Sprintf("website %s found %d jobs have golang in attribute and updated in %f sec.", provider.GetName(), length_job, time_used.Seconds()))
	}
	now := time.Now()
	u.lastTime = &now
}

func AnalysisSkill(jobDescription string) (*dao.SkillsModel, error) {
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

	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	steam := false
	req := &ollama.GenerateRequest{
		Model:  os.Getenv("OLLAMA_LLM_MODEL"), // Specify the model to use
		Stream: &steam,
		Format: json.RawMessage(`"json"`),
		Prompt: prompt,
	}

	var responseBuilder strings.Builder
	respFunc := func(resp ollama.GenerateResponse) error {
		responseBuilder.WriteString(resp.Response)
		return nil
	}

	// Use the Chat function to get a response
	ctx_10s, cancel_10s := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel_10s()
	if err := client.Generate(ctx_10s, req, respFunc); err != nil {
		return nil, err
	}

	skills := &dao.SkillsModel{}
	if err = json.Unmarshal([]byte(responseBuilder.String()), skills); err != nil {
		return nil, err
	}

	return skills, nil
}

func generateTimeToken(t time.Time) string {
	nano := t.UnixNano()

	b := make([]byte, 8)
	for i := range b {
		b[i] = byte(nano >> (i * 8))
	}

	hasher := sha256.New()
	hasher.Write(b)
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
