package clients

import (
	"time"
)

type RawData struct {
	hash         string
	linkRef      string
	jobTitle     string
	company      string
	imageCompany string
	companyInfo  string
	role         string
	subtype      JobType
	workLocation string
	position     int
	salary       string
	updatedAt    time.Time
	tags         Tags
	contact      map[Contact]string
	urgent       bool
}
