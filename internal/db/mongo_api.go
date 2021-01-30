package db

import (
	"context"

	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
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
	opts   *options.DeleteOptions
}

func New(conf *config.DatabaseConfig, ctx context.Context) (registry.Storage, error) {
	clientOptions := options.Client().ApplyURI(conf.Address)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database(conf.Database)

	mongoSt := &MongoStorage{
		client: client,
		db:     db,
		coll:   client.Database(conf.Database).Collection(conf.Collection),
		opts: options.Delete().SetCollation(&options.Collation{
			Locale:    "en_US",
			Strength:  1,
			CaseLevel: false,
		}),
	}

	return mongoSt, nil
}

func (m *MongoStorage) CreateTask(task *task.Task) (uuid.UUID, error) {
	id := uuid.New()
	task.UUID = id
	_, err := m.coll.InsertOne(context.TODO(), task)

	return id, err
}

func (m *MongoStorage) DeleteTask(id uuid.UUID) error {
	//nolint:govet
	_, err := m.coll.DeleteOne(context.TODO(), bson.D{{"uuid", id.String()}}, m.opts)
	return err
}

func (m *MongoStorage) GetTaskByUUID(id uuid.UUID) (task.Task, error) {
	var t task.Task
	//nolint:govet
	err := m.coll.FindOne(context.TODO(), bson.D{{"uuid", id.String()}}, options.FindOne()).Decode(&t)
	return t, err
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
