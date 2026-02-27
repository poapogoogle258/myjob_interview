package jobthai_test

import (
	"fmt"
	"testing"

	"github.com/poapogoogle258/myjob_interview/internel/clients/jobthai"
)

func TestFetchJobs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		keyword string
		page    int
		want    *jobthai.JobResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test fetch golang page 1",
			keyword: "golang",
			page:    1,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := 1

		FetchJob:
			result, gotErr := jobthai.FetchJobs(tt.keyword, tt.page)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FetchJobs() failed: %v", gotErr)
				}
				return
			}

			for _, job := range result.Data.SearchJobs.Data.Data {
				fmt.Printf("%d) id: %d\nตำแหน่ง: %s\nบริษัท: %s\nเงินเดือน: %s\n\n",
					i, job.ID, job.JobTitle, job.CompanyName, job.Salary)
				i++
			}

			if i < result.Data.SearchJobs.Data.Total {
				tt.page++
				goto FetchJob
			}

		})
	}
}

func TestFetchJobsDetail(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      int
		want    *jobthai.JobDetailResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "example fetch id 1609516",
			id:      1609516,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, gotErr := jobthai.FetchJobsDetail(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FetchJobsDetail() failed: %v", gotErr)
				}
				return
			}

			fmt.Println(detail.Data.GetJobRawData.Data.Description)
			// for _, propertie := range detail.Data.GetJobRawData.Data.Properties {
			// 	fmt.Printf(propertie)
			// }

		})
	}
}
