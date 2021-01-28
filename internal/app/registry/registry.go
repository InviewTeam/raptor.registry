package registry

import (
	"encoding/json"

	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/config"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

type Storage interface {
	CreateTask(*task.Task) (uuid.UUID, error)
	DeleteTask(uuid.UUID) error
	GetTasks() ([]task.Task, error)
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

func (r *Registry) CreateTask(task *task.Task) (uuid.UUID, error) {
	return r.storage.CreateTask(task)
}

func (r *Registry) DeleteTask(id uuid.UUID) error {
	return r.storage.DeleteTask(id)
}

func (r *Registry) SendTask(task *task.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.rmq.Send(data, r.conf.Rabbit.WorkerQueue)
}

func (r *Registry) GetTasks() ([]task.Task, error) {
	return r.storage.GetTasks()
}
