package format

import "github.com/google/uuid"

type TaskList []Task

type Task struct {
	UUID     uuid.UUID `json:"uuid"`
	CameraIP string    `json:"camera_ip"`
	Job      string    `json:"job"`
	Status   string    `json:"status"`
}
