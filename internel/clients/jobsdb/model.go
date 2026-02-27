package jobsdb

import "time"

type JobResponse struct {
	Results struct {
		JobIds             []string `json:"jobIds"`
		Source             string   `json:"source"`
		TotalCount         int      `json:"totalCount"`
		TotalCountNewToYou int      `json:"totalCountNewToYou"`
		TotalPages         int      `json:"totalPages"`
		NewSince           int      `json:"newSince"`
		Results            struct {
			Summary struct {
				DisplayTotalCount string `json:"displayTotalCount"`
				Text              string `json:"text"`
			} `json:"summary"`
			Jobs []struct {
				Advertiser struct {
					ID          string `json:"id"`
					Description string `json:"description"`
				} `json:"advertiser"`
				BulletPoints []string `json:"bulletPoints"`
				Branding     struct {
					SerpLogoURL string `json:"serpLogoUrl"`
				} `json:"branding"`
				CompanyName                    string    `json:"companyName,omitempty"`
				CompanyProfileStructuredDataID int       `json:"companyProfileStructuredDataId"`
				DisplayType                    string    `json:"displayType"`
				ID                             string    `json:"id"`
				IsFeatured                     bool      `json:"isFeatured"`
				ListingDate                    time.Time `json:"listingDate"`
				ListingDateDisplay             string    `json:"listingDateDisplay"`
				Locations                      []struct {
					Label        string `json:"label"`
					CountryCode  string `json:"countryCode"`
					SeoHierarchy []struct {
						ContextualName string `json:"contextualName"`
					} `json:"seoHierarchy"`
				} `json:"locations"`
				RoleID      string   `json:"roleId"`
				SalaryLabel string   `json:"salaryLabel"`
				Teaser      string   `json:"teaser"`
				Title       string   `json:"title"`
				Tracking    string   `json:"tracking"`
				WorkTypes   []string `json:"workTypes"`
			} `json:"jobs"`
		} `json:"results"`
	} `json:"results"`
}
