package registry

import (
	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

type Storage interface {
	CreateTask(*task.Task) (uuid.UUID, error)
	GetTasks() ([]task.Task, error)
}

type Registry struct {
	storage Storage
}

func New(st Storage) *Registry {
	return &Registry{
		storage: st,
	}
}

func (r *Registry) CreateTask(task *task.Task) (uuid.UUID, error) {
	return r.storage.CreateTask(task)
}

func (r *Registry) GetTasks() ([]task.Task, error) {
	return r.storage.GetTasks()
}
