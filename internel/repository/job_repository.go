package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/poapogoogle258/myjob_interview/internel/model"
)

type JobRepository interface {
	CreateIndexes(ctx context.Context) error
	UpsertByExternalID(ctx context.Context, job *model.JobModel) error
	GetByExternalID(ctx context.Context, source, externalID string) (*model.JobModel, error)
	GetPaginated(ctx context.Context, orderby string, page, limit int64) ([]*model.JobModel, error)
	GetAll(ctx context.Context) ([]*model.JobModel, error)
	GetByHashId(ctx context.Context, hashId string) (*model.JobModel, error)
}

type jobRepository struct {
	collection *mongo.Collection
}

func NewJobRepository(db *mongo.Database) JobRepository {
	return &jobRepository{
		collection: db.Collection("jobs"),
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

// GetByExternalID retrieves a single job by its source and external ID.
func (r *jobRepository) GetByExternalID(ctx context.Context, source, externalID string) (*model.JobModel, error) {
	filter := bson.M{
		"source":      source,
		"external_id": externalID,
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

// GetPaginated retrieves a list of jobs with pagination and sorting.
// orderby can be prefixed with '-' for descending order (e.g., "-posted_at").
func (r *jobRepository) GetPaginated(ctx context.Context, orderby string, page, limit int64) ([]*model.JobModel, error) {
	findOptions := options.Find()

	if orderby != "" {
		sortDir := 1
		sortKey := orderby
		if orderby[0] == '-' {
			sortDir = -1
			sortKey = orderby[1:]
		}
		findOptions.SetSort(bson.D{{Key: sortKey, Value: sortDir}})
	}

	if page < 1 {
		page = 1
	}
	skip := (page - 1) * limit
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
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

// GetAll retrieves all jobs from the collection.
func (r *jobRepository) GetAll(ctx context.Context) ([]*model.JobModel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
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
