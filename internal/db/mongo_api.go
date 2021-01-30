package db

import (
	"context"

	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
	"gitlab.com/inview-team/raptor_team/registry/pkg/format"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client         *mongo.Client
	db             *mongo.Database
	workers_coll   *mongo.Collection
	analyzers_coll *mongo.Collection
	reports_coll   *mongo.Collection
	opts           *options.DeleteOptions
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
		client:         client,
		db:             db,
		workers_coll:   client.Database(conf.Database).Collection(conf.WorkersCollection),
		analyzers_coll: client.Database(conf.Database).Collection(conf.AnalyzersCollection),
		reports_coll:   client.Database(conf.Database).Collection(conf.ReportsCollection),
		opts: options.Delete().SetCollation(&options.Collation{
			Locale:    "en_US",
			Strength:  1,
			CaseLevel: false,
		}),
	}

	return mongoSt, nil
}

func (m *MongoStorage) CreateTask(task format.Task) (uuid.UUID, error) {
	id := uuid.New()
	task.UUID = id
	_, err := m.workers_coll.InsertOne(context.TODO(), task)

	return id, err
}

func (m *MongoStorage) DeleteTask(id uuid.UUID) error {
	//nolint:govet
	_, err := m.workers_coll.DeleteOne(context.TODO(), bson.D{{"uuid", id.String()}}, m.opts)
	return err
}

func (m *MongoStorage) GetTaskByUUID(id uuid.UUID) (format.Task, error) {
	var t format.Task
	//nolint:govet
	err := m.workers_coll.FindOne(context.TODO(), bson.D{{"uuid", id.String()}}, options.FindOne()).Decode(&t)
	return t, err
}

func (m *MongoStorage) GetTasks() ([]format.Task, error) {
	var tasks []format.Task
	cur, err := m.workers_coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var t format.Task
		err := cur.Decode(&t)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tasks, err
}

func (m *MongoStorage) GetAnalyzers() ([]format.Analyzer, error) {
	var analyzers []format.Analyzer
	cur, err := m.analyzers_coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var a format.Analyzer
		err := cur.Decode(&a)
		if err != nil {
			return nil, err
		}
		analyzers = append(analyzers, a)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return analyzers, err
}

func (m *MongoStorage) CreateAnalyzer(analyzer format.Analyzer) error {
	_, err := m.analyzers_coll.InsertOne(context.TODO(), analyzer)
	return err
}

func (m *MongoStorage) GetAnalyzerByName(name string) (format.Analyzer, error) {
	var a format.Analyzer
	//nolint:govet
	err := m.analyzers_coll.FindOne(context.TODO(), bson.D{{"name", name}}, options.FindOne()).Decode(&a)
	return a, err
}

func (m *MongoStorage) DeleteAnalyzer(name string) error {
	//nolint:govet
	_, err := m.analyzers_coll.DeleteOne(context.TODO(), bson.D{{"name", name}}, m.opts)
	return err
}
