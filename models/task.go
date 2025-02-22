package models

import "time"

type Task struct {
	Header      string    `json:"header" validate:"required,min=1"`
	Description string    `json:"description" validate:"required,min=1"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
}

type TaskInfo struct {
	ID          string    `json:"task_id"`
	Header      string    `json:"header"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time,string"`
	EndTime     time.Time `json:"end_time,string"`
}
