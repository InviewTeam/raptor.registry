package db

import "gitlab.com/inview-team/raptor_team/registry/task"

type Storage interface {
	AddTask(*task.Task) error
	GetTasks() ([]*task.Task, error)
}
