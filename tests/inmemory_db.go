package tests

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/task"
)

type InMemoryDB struct {
	tasks []task.Task
}

func NewDB() registry.Storage {
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

func (im *InMemoryDB) GetTaskByUUID(id uuid.UUID) (task.Task, error) {
	for i := range im.tasks {
		if im.tasks[i].UUID == id {
			return im.tasks[i], nil
		}
	}
	return task.Task{}, fmt.Errorf("no task with id %s", id.String())
}

func (im *InMemoryDB) DeleteTask(id uuid.UUID) error {
	for i := range im.tasks {
		if im.tasks[i].UUID == id {
			im.tasks[i] = im.tasks[len(im.tasks)-1]
			im.tasks = im.tasks[:len(im.tasks)-1]
			return nil
		}
	}
	return fmt.Errorf("no task with id %s", id.String())
}
