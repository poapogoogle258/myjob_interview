package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FetchJobs(keyword string, page int) (*Jobthai_jobResponse, error) {
	sha256Hash := "8c21badbcb9da924a3ed99c6d2f16d34758a045523495b4458f2a970c70cd0b2"
	variables := fmt.Sprintf(`{"searchJobsFilter":{"province":"01","keyword":"%s","l":"th","page":%d},"orderBy":"UPDATED_AT_DESC","staticDataVersion":{"jobType":null,"subjobType":null}}`, keyword, page)
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

	var jobData Jobthai_jobResponse
	if err := json.Unmarshal(body, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil
}

func FetchJobsDetail(id int) (*Jobthai_jobDetailResponse, error) {
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

	var jobData Jobthai_jobDetailResponse
	if err := json.Unmarshal(body, &jobData); err != nil {
		return nil, err
	}

	return &jobData, nil

}
