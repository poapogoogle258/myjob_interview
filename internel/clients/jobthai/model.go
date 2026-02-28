package jobthai

import (
	"time"
)

type JobResponse struct {
	Data struct {
		SearchJobs struct {
			Data struct {
				Total int `json:"total"`
				Data  []struct {
					ID          int    `json:"id"`
					CompanyID   int    `json:"companyID"`
					JobTitle    string `json:"jobTitle"`
					CompanyName string `json:"companyName"`
					CompanyLogo string `json:"companyLogo"`
					Province    struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"province"`
					DisabledPerson struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					} `json:"disabledPerson"`
					Country struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					} `json:"country"`
					WorkLocation string `json:"workLocation"`
					Salary       string `json:"salary"`
					Urgent       struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					} `json:"urgent"`
					JobType struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					} `json:"jobType"`
					Tags      []string  `json:"tags"`
					UpdatedAt time.Time `json:"updatedAt"`
				} `json:"data"`
			} `json:"data"`
		} `json:"searchJobs"`
	} `json:"data"`
}

type JobDetailResponse struct {
	Data struct {
		GetJobRawData struct {
			Data struct {
				ID      int    `json:"_id"`
				Title   string `json:"title"`
				Company struct {
					ID             int      `json:"_id"`
					Name           string   `json:"name"`
					Logo           string   `json:"logo"`
					Website        string   `json:"website"`
					Pictures       []string `json:"pictures"`
					Detail         string   `json:"detail"`
					Benefit        string   `json:"benefit"`
					SpecialContent string   `json:"specialContent"`
				} `json:"company"`
				Properties   []string `json:"properties"`
				Benefit      string   `json:"benefit"`
				ApplyMethod  string   `json:"applyMethod"`
				Description  string   `json:"description"`
				WorkLocation struct {
					Address  string `json:"address"`
					Country  string `json:"country"`
					Province struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"province"`
					District struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"district"`
					Subdistrict struct {
						Name string `json:"name"`
					} `json:"subdistrict"`
					Industrial     string `json:"industrial"`
					IndustrialName string `json:"industrialName"`
					Map            string `json:"map"`
					Direction      string `json:"direction"`
				} `json:"workLocation"`
				NumberOfPosition string `json:"numberOfPosition"`
				Contact          struct {
					Name     string `json:"name"`
					Tel      string `json:"tel"`
					Location struct {
						Address  string `json:"address"`
						Province struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"province"`
						District struct {
							Name string `json:"name"`
						} `json:"district"`
						Subdistrict struct {
							Name string `json:"name"`
						} `json:"subdistrict"`
						IndustrialName string `json:"industrialName"`
						Country        struct {
							Name string `json:"name"`
						} `json:"country"`
						Zipcode string `json:"zipcode"`
					} `json:"location"`
					LineID string   `json:"lineID"`
					Fax    string   `json:"fax"`
					Emails []string `json:"emails"`
				} `json:"contact"`
				DisabledPerson struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"disabledPerson"`
				EnglishApply bool   `json:"englishApply"`
				Salary       string `json:"salary"`
				Urgent       struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"urgent"`
				TransitStations []string `json:"transitStations"`
				JobType         struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"jobType"`
				SubjobType struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"subjobType"`
				UpdatedAt         time.Time `json:"updatedAt"`
				Website           string    `json:"website"`
				ApplyExternalLink string    `json:"applyExternalLink"`
				Tags              []string  `json:"tags"`
				EmploymentType    string    `json:"employmentType"`
			} `json:"data"`
		} `json:"getJobRawData"`
	} `json:"data"`
}
