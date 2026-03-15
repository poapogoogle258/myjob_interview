package jobsdb_test

import (
	"fmt"
	"testing"

	"github.com/poapogoogle258/myjob_interview/internel/client/jobsdb"
)

func TestFetchJob(t *testing.T) {
	i, page := 1, 1
FetchData:
	data, _ := jobsdb.FetchJob(page)
	for _, job := range data.Results.Results.Jobs {
		fmt.Println(i, ") ", job.Title)
		i++
	}

	if len(data.Results.Results.Jobs) != 0 {
		page++
		goto FetchData
	}

}

func TestFetchJobDetail(t *testing.T) {
	data, err := jobsdb.FetchJobDetail("90944662")
	if err != nil {
		t.Errorf("FetchJobDetail() failed: %v", err)
		return
	}

	fmt.Println(data.Jobdetails)

}

func TestClient_FetchJobs(t *testing.T) {
	var c jobsdb.Client
	result, err := c.FetchJobs()
	if err != nil {
		t.Errorf("FetchJobs() failed: %v", err)
		return
	}

	for job := range result {
		fmt.Println(job)
	}

}
