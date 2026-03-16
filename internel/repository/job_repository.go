package repository

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	model "github.com/poapogoogle258/myjob_interview/internel/model/dao"
)

type JobRepository interface {
	GetAll(ctx context.Context) ([]*model.JobModel, error)
	GetByHashId(ctx context.Context, hashId string) (*model.JobModel, error)
	GetExceptSyncId(ctx context.Context, source string, syncId string) ([]*model.JobModel, error)
	CreateIndexes(ctx context.Context) error
	UpsertByExternalID(ctx context.Context, job *model.JobModel) error
	UpdateStatus(ctx context.Context, jobID string, status string) error
	IsExist(ctx context.Context, jobID string) bool
}

type jobRepository struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewJobRepository(db *mongo.Database, logger *slog.Logger) JobRepository {
	return &jobRepository{
		collection: db.Collection("jobs"),
		logger:     logger,
	}
}

// CreateIndexes ensures that the collection has the necessary indexes for performance and uniqueness.
func (r *jobRepository) CreateIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "source", Value: 1},
			{Key: "external_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// UpsertByExternalID saves a job to the database.
// If a job with the same Source and ExternalID exists, it updates the existing record.
func (r *jobRepository) UpsertByExternalID(ctx context.Context, job *model.JobModel) error {
	now := time.Now()

	filter := bson.M{
		"external_id": job.ExternalID,
		"hash_id":     job.HashId,
		"source":      job.Source,
	}

	job.UpdatedAt = now
	// We set CreatedAt to zero so it's omitted from the $set operation due to omitempty.
	// This prevents overwriting the original creation time on updates.
	job.CreatedAt = time.Time{}

	update := bson.M{
		"$set":         job,
		"$setOnInsert": bson.M{"created_at": now},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *jobRepository) GetByHashId(ctx context.Context, hashId string) (*model.JobModel, error) {
	filter := bson.M{
		"hash_id": hashId,
	}

	var job model.JobModel
	err := r.collection.FindOne(ctx, filter).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// GetAll retrieves all jobs from the collection.
func (r *jobRepository) GetAll(ctx context.Context) ([]*model.JobModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": bson.M{"$ne": "deleted"}, "skills.languages": bson.M{"$in": []string{"golang", "go"}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*model.JobModel
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepository) UpdateStatus(ctx context.Context, jobID string, status string) error {

	objID, _ := primitive.ObjectIDFromHex(jobID)
	statusLog := bson.M{
		"status":     status,
		"updated_at": time.Now(),
	}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
		"$push": bson.M{
			"status_log": statusLog,
		},
	}

	r.collection.UpdateOne(ctx, bson.M{"_id": objID, "status_log": bson.M{"$not": bson.M{"$size": 0}}}, bson.M{"$set": bson.M{"status_log": bson.A{}}})
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)

	return err
}

func (r *jobRepository) IsExist(ctx context.Context, jobID string) bool {
	objID, _ := primitive.ObjectIDFromHex(jobID)
	filter := bson.M{"_id": objID}
	result := model.JobModel{}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	return err != mongo.ErrNoDocuments
}

func (r *jobRepository) GetExceptSyncId(ctx context.Context, source string, syncId string) ([]*model.JobModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"source": source, "sync_id": bson.M{"$ne": syncId}, "status": bson.M{"$in": []string{
		"new",
		"favorite",
		"viewed",
		"registered",
		"optional",
	}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []*model.JobModel
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}
