package tests

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/inview-team/raptor_team/registry/internal/app/registry"
	"gitlab.com/inview-team/raptor_team/registry/pkg/format"
)

type InMemoryDB struct {
	tasks []format.Task
}

func NewDB() registry.Storage {
	return &InMemoryDB{
		tasks: []format.Task{},
	}
}

func (im *InMemoryDB) CreateTask(task format.Task) (uuid.UUID, error) {
	id := uuid.New()
	task.UUID = id
	im.tasks = append(im.tasks, task)
	return id, nil
}

func (im *InMemoryDB) GetTasks() ([]format.Task, error) {
	return im.tasks, nil
}

func (im *InMemoryDB) GetTaskByUUID(id uuid.UUID) (format.Task, error) {
	for i := range im.tasks {
		if im.tasks[i].UUID == id {
			return im.tasks[i], nil
		}
	}
	return format.Task{}, fmt.Errorf("no task with id %s", id.String())
}

func (im *InMemoryDB) UpdateTask(id uuid.UUID, key, value string) error {
	if key != "status" {
		return fmt.Errorf("unknown key %s", key)
	}
	for i := range im.tasks {
		if im.tasks[i].UUID == id {
			im.tasks[i].Status = value
			return nil
		}
	}
	return fmt.Errorf("no task with id %s", id.String())
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

func (im *InMemoryDB) CreateAnalyzer(format.Analyzer) error {
	//TODO: analyzer tests
	return nil
}

func (im *InMemoryDB) DeleteAnalyzer(string) error {
	//TODO: analyzer tests
	return nil
}

func (im *InMemoryDB) GetAnalyzerByName(string) (format.Analyzer, error) {
	//TODO: analyzer tests
	return format.Analyzer{}, nil
}

func (im *InMemoryDB) GetAnalyzers() ([]format.Analyzer, error) {
	//TODO: analyzer tests
	return []format.Analyzer{}, nil
}

func (im *InMemoryDB) AddReport(format.Report) error {
	return nil
}

func (im *InMemoryDB) GetReport(uuid.UUID) (format.Report, error) {
	return format.Report{}, nil
}
