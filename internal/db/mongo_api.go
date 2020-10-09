package db

import (
	"context"
	"fmt"

	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	//_ "gitlab.com/inview-team/raptor_team/registry/internal/db/migrations" // database migrations
	"gitlab.com/inview-team/raptor_team/registry/task"

	"github.com/google/uuid"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	db *mongo.Database
}

func InitDB(host, user, password, database string, ctx context.Context) (*mongo.Database, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:27017", user, password, host)
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
	migrate.SetDatabase(db)
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		return nil, err
	}
	return db, nil
}

func New(db *mongo.Database) registry.Storage {
	return &MongoStorage{db: db}
}

func (m *MongoStorage) CreateTask(task *task.Task) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *MongoStorage) GetTasks() ([]task.Task, error) {
	return nil, nil
}
