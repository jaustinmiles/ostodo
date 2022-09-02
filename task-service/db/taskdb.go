package db

import (
	"github.com/jaustinmiles/ostodo/task-service/common"
	"github.com/jaustinmiles/ostodo/task-service/tasks"
)

func Run() {
	l := common.GetLogger()
	db := NewTaskDatabase()
	if err := db.Connect(); err != nil {
		l.Errorf("failed to connect to database")
		return
	}

	task := tasks.GetTestTask()
	if err := db.InsertTask(task); err != nil {
		l.Error("error inserting into db, shutting down", err)
		return
	}

	db.Close()
}
