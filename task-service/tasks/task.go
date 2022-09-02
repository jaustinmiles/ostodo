package tasks

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	User           string
	UUID           uuid.UUID
	Name           string
	CompletionTime time.Time
	Repetitions    int
	Dependencies   []string
}

func GetTestTask() Task {
	return Task{
		User:           "jaustinmiles",
		UUID:           uuid.New(),
		Name:           "example task",
		CompletionTime: time.Now().Add(time.Hour),
		Repetitions:    0,
		Dependencies:   nil,
	}
}
