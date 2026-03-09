package jobthai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/poapogoogle258/myjob_interview/internel/model"
	provider "github.com/poapogoogle258/myjob_interview/internel/provider/clients_provider"
)

type Client struct{}

func (c *Client) GetName() string {
	return "jobthai"
}

func (c *Client) FetchJobs() ([]*model.JobModel, error) {
	results_nums := make([]*model.JobModel, 0)
	map_results := make(map[string]struct{})
	keywords := []string{"go", "golang"}

	for _, keyword := range keywords {
		results := make([]*model.JobModel, 0, 100)
		page := 1
		for {
			jobs, err := FetchJobs(keyword, page)
			if err != nil {
				return nil, err
			}
			if len(jobs.Data.SearchJobs.Data.Data) == 0 {
				break
			}
			for _, job := range jobs.Data.SearchJobs.Data.Data {
				res, err := FetchJobsDetail(job.ID)
				if err != nil {
					return nil, err
				}
				// convert response to Job DAO
				detail := res.Data.GetJobRawData.Data

				// check duplicate jobs fetched
				if _, ok := map_results[fmt.Sprintf("%d", job.ID)]; ok {
					continue
				}

				results = append(results, &model.JobModel{
					Source:      c.GetName(),
					ExternalID:  fmt.Sprintf("%d", job.ID),
					Title:       detail.Title,
					CompanyName: detail.Company.Name,
					Location:    fmt.Sprintf("%s %s %s %s", detail.WorkLocation.Province.Name, detail.WorkLocation.District.Name, detail.WorkLocation.Subdistrict.Name, detail.WorkLocation.Address),
					Salary:      detail.Salary,
					Description: fmt.Sprintf("properties :\n%s\n description :\n %s", detail.Description, strings.Join(detail.Properties, "\n")),
					Status:      "new",
					Skills:      nil,
					URL:         fmt.Sprintf("https://www.jobthai.com/th/job/%d", job.ID),
					PostedAt:    detail.UpdatedAt,
				})
				map_results[fmt.Sprintf("%d", job.ID)] = struct{}{}
			}
			page++
		}

		results_nums = append(results_nums, results...)
	}

	return results_nums, nil
}

func init() {
	provider.Register(&Client{})
}

func FetchJobs(keyword string, page int) (*JobResponse, error) {
	sha256Hash := "8c21badbcb9da924a3ed99c6d2f16d34758a045523495b4458f2a970c70cd0b2"
	variables := fmt.Sprintf(`{"searchJobsFilter":{"region":"6","jobtype":"7","keyword":"%s","l":"th","page":%d},"orderBy":"UPDATED_AT_DESC","staticDataVersion":{"jobType":null,"subjobType":null}}`, keyword, page)
	extensions := fmt.Sprintf(`{"persistedQuery":{"version":1,"sha256Hash":"%s"}}`, sha256Hash)

	u, _ := url.Parse("https://api.jobthai.com/v1/graphql")
	q := u.Query()
	q.Set("variables", variables)
	q.Set("extensions", extensions)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("client-name", "jobthai-upgrade-mobile")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5)")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apollo-require-preflight", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jobData JobResponse
	if err := json.Unmarshal(body, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil
}

func FetchJobsDetail(id int) (*JobDetailResponse, error) {
	sha256Hash := "4fe6bb56592bed522754f07a0cf519150f67ad761c375917b9084f216f0ea93e"
	variables := fmt.Sprintf(`{"id":%d,"l":"th","isJobbuffer":false,"staticDataVersion":{"jobType":null,"subjobType":null}}`, id)
	extensions := fmt.Sprintf(`{"persistedQuery":{"version":1,"sha256Hash":"%s"}}`, sha256Hash)

	u, _ := url.Parse("https://api.jobthai.com/v1/graphql")
	q := u.Query()
	q.Set("variables", variables)
	q.Set("extensions", extensions)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("client-name", "jobthai-upgrade-mobile")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5)")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apollo-require-preflight", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jobData JobDetailResponse
	if err := json.Unmarshal(body, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil
}
