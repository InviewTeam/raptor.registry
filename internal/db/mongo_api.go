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
}

func New(conf *config.Settings, ctx context.Context) (registry.Storage, error) {
	clientOptions := options.Client().ApplyURI(conf.DatabaseAddress)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	db := client.Database(conf.DatabaseName)

	mongoSt := &MongoStorage{
		client:         client,
		db:             db,
		workers_coll:   client.Database(conf.DatabaseName).Collection(conf.WorkersCollection),
		analyzers_coll: client.Database(conf.DatabaseName).Collection(conf.AnalyzersCollection),
		reports_coll:   client.Database(conf.DatabaseName).Collection(conf.ReportsCollection),
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
	_, err := m.workers_coll.DeleteOne(context.TODO(), bson.M{"uuid": id.String()}, options.Delete())
	return err
}

func (m *MongoStorage) UpdateTask(id uuid.UUID, key, value string) error {
	_, err := m.workers_coll.UpdateOne(context.TODO(), bson.M{"uuid": id}, bson.M{key: value}, options.Update().SetUpsert(false))
	return err
}

func (m *MongoStorage) GetTaskByUUID(id uuid.UUID) (format.Task, error) {
	var t format.Task
	err := m.workers_coll.FindOne(context.TODO(), bson.M{"uuid": id.String()}, options.FindOne()).Decode(&t)
	return t, err
}

func (m *MongoStorage) GetTasks() ([]format.Task, error) {
	var tasks []format.Task
	cur, err := m.workers_coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	err = cur.All(context.Background(), &tasks)
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
	err := m.analyzers_coll.FindOne(context.TODO(), bson.M{"name": name}, options.FindOne()).Decode(&a)
	return a, err
}

func (m *MongoStorage) DeleteAnalyzer(name string) error {
	_, err := m.analyzers_coll.DeleteOne(context.TODO(), bson.M{"name": name}, options.Delete())
	return err
}

func (m *MongoStorage) AddReport(rep format.Report) error {
	_, err := m.reports_coll.InsertOne(context.TODO(), rep)
	return err
}

func (m *MongoStorage) GetReport(id uuid.UUID) (format.Report, error) {
	var rep format.Report
	err := m.reports_coll.FindOne(context.TODO(), bson.M{"uuid": id.String()}, options.FindOne()).Decode(&rep)
	return rep, err
}
