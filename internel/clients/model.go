package clients

import (
	"time"
)

type RawData struct {
	Id             string
	LinkRef        string
	JobTitle       string
	JobDescription string
	Company        string
	CompanyInfo    string
	Subtype        JobType
	WorkLocation   string
	Salary         string
	UpdatedAt      time.Time
	Tags           Tags
	Contact        map[Contact]string
	Urgent         bool
}
