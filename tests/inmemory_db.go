package tests

import (
	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

type InMemoryDB struct {
	tasks []task.Task
}

func New() registry.Storage {
	return &InMemoryDB{
		tasks: []task.Task{},
	}
}

func (im *InMemoryDB) CreateTask(task *task.Task) (uuid.UUID, error) {
	id := uuid.New()
	task.UUID = id
	im.tasks = append(im.tasks, *task)
	return id, nil
}

func (im *InMemoryDB) GetTasks() ([]task.Task, error) {
	return im.tasks, nil
}
