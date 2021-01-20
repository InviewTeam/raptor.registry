package registry

import (
	"encoding/json"

	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

type Storage interface {
	CreateTask(*task.Task) (uuid.UUID, error)
	GetTasks() ([]task.Task, error)
}

type PublisherInterface interface {
	Connect() error
	Close() error
	Send([]byte) error
}

type Registry struct {
	storage Storage
	rmq     PublisherInterface
}

func New(st Storage, pub PublisherInterface) *Registry {
	return &Registry{
		storage: st,
		rmq:     pub,
	}
}

func (r *Registry) CreateTask(task *task.Task) (uuid.UUID, error) {
	return r.storage.CreateTask(task)
}

func (r *Registry) SendTask(task *task.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.rmq.Send(data)
}

func (r *Registry) GetTasks() ([]task.Task, error) {
	return r.storage.GetTasks()
}
