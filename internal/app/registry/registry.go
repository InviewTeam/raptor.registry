package registry

import (
	"encoding/json"

	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
	"gitlab.com/inview-team/raptor_team/registry/pkg/format"
)

type Storage interface {
	CreateTask(format.Task) (uuid.UUID, error)
	DeleteTask(uuid.UUID) error
	GetTaskByUUID(uuid.UUID) (format.Task, error)
	GetTasks() ([]format.Task, error)

	CreateAnalyzer(format.Analyzer) error
	DeleteAnalyzer(string) error
	GetAnalyzerByName(string) (format.Analyzer, error)
	GetAnalyzers() ([]format.Analyzer, error)
}

type PublisherInterface interface {
	Connect() error
	Close() error
	DeclareQueue(string) error
	Send([]byte, string) error
}

type Registry struct {
	storage Storage
	rmq     PublisherInterface
	conf    *config.Settings
}

func New(conf *config.Settings, st Storage, pub PublisherInterface) *Registry {
	return &Registry{
		storage: st,
		rmq:     pub,
		conf:    conf,
	}
}

func (r *Registry) CreateTask(task format.Task) (uuid.UUID, error) {
	return r.storage.CreateTask(task)
}

func (r *Registry) DeleteTask(id uuid.UUID) error {
	return r.storage.DeleteTask(id)
}

func (r *Registry) SendTask(task *format.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.rmq.Send(data, r.conf.Rabbit.WorkerQueue)
}

func (r *Registry) GetTaskByUUID(id uuid.UUID) (format.Task, error) {
	return r.storage.GetTaskByUUID(id)
}

func (r *Registry) GetTasks() ([]format.Task, error) {
	return r.storage.GetTasks()
}

func (r *Registry) StopTask(id uuid.UUID) error {
	req, err := json.Marshal(map[string]string{"uuid": id.String(), "status": "done"})
	if err != nil {
		return err
	}
	return r.rmq.Send(req, r.conf.Rabbit.WorkerQueue)
}

func (r *Registry) GetAnalyzers() ([]format.Analyzer, error) {
	return r.storage.GetAnalyzers()
}

func (r *Registry) GetAnalyzerByName(name string) (format.Analyzer, error) {
	return r.storage.GetAnalyzerByName(name)
}

func (r *Registry) CreateAnalyzer(analyzer format.Analyzer) error {
	return r.storage.CreateAnalyzer(analyzer)
}

func (r *Registry) DeleteAnalyzer(name string) error {
	return r.storage.DeleteAnalyzer(name)
}
