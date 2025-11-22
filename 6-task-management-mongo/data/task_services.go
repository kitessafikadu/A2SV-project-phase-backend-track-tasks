package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"task-management/models"
)

var (
	ErrNotFound    = errors.New("task not found")
	ErrInvalidDate = errors.New("invalid date format, use RFC3339")
)

var taskCollection *mongo.Collection

func InitMongo(client *mongo.Client) {
	taskCollection = client.Database("taskdb").Collection("tasks")
}

func CreateTask(input models.TaskInput) (models.Task, error) {
	due := time.Time{}
	var err error
	if input.DueDate != "" {
		due, err = time.Parse(time.RFC3339, input.DueDate)
		if err != nil {
			return models.Task{}, ErrInvalidDate
		}
	}

	t := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Description: input.Description,
		DueDate:     due,
		Status:      input.Status,
	}

	_, err = taskCollection.InsertOne(context.Background(), t)
	if err != nil {
		return models.Task{}, err
	}

	return t, nil
}

func GetAllTasks() ([]models.Task, error) {
	cur, err := taskCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var tasks []models.Task
	for cur.Next(context.Background()) {
		var t models.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTask(id string) (models.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, ErrNotFound
	}

	var task models.Task
	err = taskCollection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return models.Task{}, ErrNotFound
	}

	return task, nil
}

func UpdateTask(id string, input models.TaskInput) (models.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, ErrNotFound
	}

	due := time.Time{}
	if input.DueDate != "" {
		due, err = time.Parse(time.RFC3339, input.DueDate)
		if err != nil {
			return models.Task{}, ErrInvalidDate
		}
	}

	update := bson.M{
		"$set": bson.M{
			"title":       input.Title,
			"description": input.Description,
			"status":      input.Status,
			"due_date":    due,
		},
	}

	_, err = taskCollection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return models.Task{}, ErrNotFound
	}

	return GetTask(id)
}

func DeleteTask(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}

	res, err := taskCollection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
