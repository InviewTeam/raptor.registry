package db

import (
	"context"
	"fmt"

	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/task"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

func New(host, user, password, database, coll string, ctx context.Context) (registry.Storage, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s", user, password, host)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database(database)

	mongoSt := &MongoStorage{
		client: client,
		db:     db,
		coll:   client.Database(database).Collection(coll),
	}

	return mongoSt, nil
}

func (m *MongoStorage) CreateTask(task *task.Task) (uuid.UUID, error) {
	id := uuid.New()
	task.UUID = id
	_, err := m.coll.InsertOne(context.TODO(), task)

	return id, err
}

func (m *MongoStorage) GetTasks() ([]task.Task, error) {
	var tasks []task.Task
	cur, err := m.coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var t task.Task
		err := cur.Decode(&t)
		tasks = append(tasks, t)
		if err != nil {
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tasks, err
}
