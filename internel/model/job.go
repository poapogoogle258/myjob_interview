package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobModel represents the job document structure in MongoDB
type JobModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Source      string             `bson:"source"`      // e.g., "jobsdb", "jobthai"
	ExternalID  string             `bson:"external_id"` // ID from the original platform
	Title       string             `bson:"title"`
	CompanyName string             `bson:"company_name"`
	Location    string             `bson:"location"`
	Salary      string             `bson:"salary"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
	Skills      *SkillsModel       `bson:"skills"`
	URL         string             `bson:"url"`
	PostedAt    time.Time          `bson:"posted_at"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type SkillsModel struct {
	Languages  []string `bson:"languages"`
	Frameworks []string `bson:"frameworks"`
	Tools      []string `bson:"tools"`
	Databases  []string `bson:"databases"`
	HardSkills []string `bson:"hard_skills"`
	SoftSkills []string `bson:"soft_skills"`
}
