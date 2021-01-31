package registry

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
	"gitlab.com/inview-team/raptor_team/registry/pkg/format"
)

type Storage interface {
	CreateTask(format.Task) (uuid.UUID, error)
	DeleteTask(uuid.UUID) error
	UpdateTask(uuid.UUID, string, string) error
	GetTaskByUUID(uuid.UUID) (format.Task, error)
	GetTasks() ([]format.Task, error)

	CreateAnalyzer(format.Analyzer) error
	DeleteAnalyzer(string) error
	GetAnalyzerByName(string) (format.Analyzer, error)
	GetAnalyzers() ([]format.Analyzer, error)

	AddReport(format.Report) error
	GetReport(uuid.UUID) (format.Report, error)
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
	task.Status = "in work"
	id, err := r.storage.CreateTask(task)
	if err != nil {
		return id, err
	}
	data, err := json.Marshal(task)
	if err != nil {
		return id, err

	}
	return id, r.rmq.Send(data, r.conf.Rabbit.WorkerQueue)
}

func (r *Registry) DeleteTask(id uuid.UUID) error {
	err := r.storage.DeleteTask(id)
	if err != nil {
		log.Printf("failed to delete task: %s", err.Error())
	}
	req, err := json.Marshal(map[string]string{"uuid": id.String(), "status": "stopped"})
	if err != nil {
		return err
	}
	return r.rmq.Send(req, r.conf.Rabbit.WorkerQueue)
}

func (r *Registry) GetTaskByUUID(id uuid.UUID) (format.Task, error) {
	return r.storage.GetTaskByUUID(id)
}

func (r *Registry) GetTasks() ([]format.Task, error) {
	return r.storage.GetTasks()
}

func (r *Registry) StopTask(id uuid.UUID) error {
	err := r.storage.UpdateTask(id, "status", "stopped")
	if err != nil {
		log.Printf("failed to change task status: %s", err.Error())
	}
	req, err := json.Marshal(map[string]string{"uuid": id.String(), "status": "stopped"})
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

func (r *Registry) AddReport(rep format.Report) error {
	return r.storage.AddReport(rep)
}

func (r *Registry) GetReport(id uuid.UUID) (format.Report, error) {
	return r.storage.GetReport(id)
}
