package database

import (
	"context"
	"log"
	"time"

	"github.com/demola234/golang-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString = "mongodb+srv://dbAdemola:Ademola$123@cluster0.4znphzp.mongodb.net/?retryWrites=true&w=majority"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return &DB{client: client}
}

func (db *DB) OpenCollection() *mongo.Client {
	return db.client
}

func (db *DB) GetJob(id string) *model.JobListing {
	jobCollection := db.client.Database("job-listing").Collection("jobs")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	var jobListing model.JobListing

	err := jobCollection.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollection := db.client.Database("job-listing").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var jobListing []*model.JobListing
	cusor, err := jobCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cusor.All(context.TODO(), &jobListing); err != nil {
		log.Fatal(err)
	}

	return jobListing
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {
	jobCollection := db.client.Database("job-listing").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	inserted, err := jobCollection.InsertOne(ctx, bson.M{
		"title":       jobInfo.Title,
		"description": jobInfo.Description,
		"company":     jobInfo.Company,
		"url":         jobInfo.URL,
		
	})

	if err != nil {
		log.Fatal(err)
	}

	InsertedID := inserted.InsertedID.(primitive.ObjectID).Hex()

	returningJob := model.JobListing{
		ID:          InsertedID,
		Title:       jobInfo.Title,
		Description: jobInfo.Description,
		Company:     jobInfo.Company,
	}

	return &returningJob
}

func (db *DB) UpdateJobListing(jobId string, jobModel *model.UpdateJobListingInput) *model.JobListing {
	jobCollection := db.client.Database("job-listing").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateInfo := bson.M{}

	if jobModel.Title != nil {
		updateInfo["title"] = jobModel.Title
	}
	if jobModel.Description != nil {
		updateInfo["description"] = jobModel.Description
	}
	if jobModel.Company != nil {
		updateInfo["company"] = jobModel.Company
	}
	if jobModel.URL != nil {
		updateInfo["url"] = jobModel.URL
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateInfo}

	result := jobCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate())

	var jobListing model.JobListing

	if err := result.Decode(&jobListing); err != nil {
		log.Fatal(err)
	}

	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollection := db.client.Database("job-listing").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}

	_, err := jobCollection.DeleteOne(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}

	return &model.DeleteJobResponse{DeleteJobID: jobId}
}
