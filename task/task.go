package task

import "github.com/gofrs/uuid"

type TaskList []Task

type Task struct {
	UUID     uuid.UUID `json:"uuid"`
	CameraIP string    `json:"camera_ip"`
	Jobs     []Job     `json:"jobs"`
}

type Job struct {
	Title string `json:"title"`
}
