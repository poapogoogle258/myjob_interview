package dao

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobModel represents the job document structure in MongoDB
type JobModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Source      string             `bson:"source"`      // e.g., "jobsdb", "jobthai"
	ExternalID  string             `bson:"external_id"` // ID from the original platform
	SyncID      string             `bson:"sync_id"`
	Title       string             `bson:"title"`
	CompanyName string             `bson:"company_name"`
	Location    string             `bson:"location"`
	Salary      string             `bson:"salary"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
	StatusLOG   []StatusLog        `bson:"status_log"`
	Skills      *SkillsModel       `bson:"skills"`
	URL         string             `bson:"url"`
	HashId      string             `bson:"hash_id"`
	PostedAt    time.Time          `bson:"posted_at"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type StatusLog struct {
	Status  string    `bson:"status"`
	Changed time.Time `bson:"changed_at"`
}

type SkillsModel struct {
	Languages  []string `bson:"languages"`
	Frameworks []string `bson:"frameworks"`
	Tools      []string `bson:"tools"`
	Databases  []string `bson:"databases"`
	HardSkills []string `bson:"hard_skills"`
	SoftSkills []string `bson:"soft_skills"`
}

func GetHashJobId(data *JobModel) string {
	title := strings.ToLower(strings.ReplaceAll(data.Title, " ", ""))
	company := strings.ToLower(strings.ReplaceAll(data.CompanyName, " ", ""))
	key := fmt.Sprintf("%s;%s;%s", strings.ToLower(data.Source), title, company)
	hashBytes := sha256.Sum256([]byte(key))

	return hex.EncodeToString(hashBytes[:])
}
