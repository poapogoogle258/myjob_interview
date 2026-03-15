package jobsdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/poapogoogle258/myjob_interview/internel/client/provider"
	model "github.com/poapogoogle258/myjob_interview/internel/model/dao"
)

type Client struct{}

func (c *Client) GetName() string {
	return "jobsdb"
}

func (c *Client) FetchJobs() ([]*model.JobModel, error) {
	results := make([]*model.JobModel, 0, 100)
	page := 1
	for {
		jobs, err := FetchJob(page)
		if err != nil {
			return nil, err
		}
		for _, job := range jobs.Results.Results.Jobs {
			res, err := FetchJobDetail(job.ID)
			if err != nil {
				return nil, err
			}
			// convert response to Job DAO
			detail := res.Jobdetails.Result.Job
			salary := "ไม่ได้ระบุ"
			if detail.Salary != nil && detail.Salary.Label != "" {
				salary = detail.Salary.Label
			}
			results = append(results, &model.JobModel{
				Source:      c.GetName(),
				ExternalID:  fmt.Sprintf("%s", job.ID),
				Title:       detail.Title,
				CompanyName: detail.Advertiser.Name,
				Location:    detail.Location.Label,
				Salary:      salary,
				Description: fmt.Sprintf("abstract :\n%s\n description :\n %s", detail.Abstract, detail.Content),
				Status:      "new",
				Skills:      nil,
				URL:         fmt.Sprintf("https://th.jobsdb.com/th/job/%s", job.ID),
				PostedAt:    detail.ListedAt.DateTimeUtc,
			})
		}
		if len(results) == jobs.Results.TotalCount {
			break
		}
		page++
	}

	return results, nil
}

func init() {
	provider.Register(&Client{})
}

func FetchJob(page int) (*JobResponse, error) {
	//https://th.jobsdb.com/th/golang-jobs/in-กรุงเทพมหานคร?page=1
	url := fmt.Sprintf(`https://th.jobsdb.com/th/golang-jobs/in-กรุงเทพมหานคร?page=%d`, page)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูล: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	startKey := "window.SEEK_REDUX_DATA = "
	_, after, ok := strings.Cut(html, startKey)
	if !ok {
		return nil, fmt.Errorf("ไม่เจอข้อมูล SEEK_REDUX_DATA:")
	}
	before, _, ok := strings.Cut(after, "};")
	if !ok {
		return nil, fmt.Errorf("ไม่เจอข้อมูลจบไม่เจอ")
	}

	jsonString := []byte(before + "}")

	var jobData JobResponse
	if err := json.Unmarshal(jsonString, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil
}

func FetchJobDetail(id string) (*JobDetailResponse, error) {
	// https://th.jobsdb.com/th/job/
	url := fmt.Sprintf(`https://th.jobsdb.com/th/job/%s`, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูล: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	startKey := "window.SEEK_REDUX_DATA = "
	_, after, ok := strings.Cut(html, startKey)
	if !ok {
		return nil, fmt.Errorf("ไม่เจอข้อมูล SEEK_REDUX_DATA:")
	}
	before, _, ok := strings.Cut(after, "};")
	if !ok {
		return nil, fmt.Errorf("ไม่เจอข้อมูลจบไม่เจอ")
	}

	jsonString := []byte(before + "}")

	var jobData JobDetailResponse
	if err := json.Unmarshal(jsonString, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil
}
